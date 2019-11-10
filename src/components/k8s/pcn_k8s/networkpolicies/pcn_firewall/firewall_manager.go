package pcnfirewall

import (
	"strings"
	"sync"
	"time"

	pcn_controllers "github.com/polycube-network/polycube/src/components/k8s/pcn_k8s/controllers"
	pcn_types "github.com/polycube-network/polycube/src/components/k8s/pcn_k8s/types"
	"github.com/polycube-network/polycube/src/components/k8s/utils"
	k8sfirewall "github.com/polycube-network/polycube/src/components/k8s/utils/k8sfirewall"
	core_v1 "k8s.io/api/core/v1"
	k8s_types "k8s.io/apimachinery/pkg/types"
)

// PcnFirewallManager is the interface of the firewall manager.
type PcnFirewallManager interface {
	// Link adds a new pod to the list of pods that must be managed by this manager.
	// Best practice is to only link similar pods
	// (i.e.: same labels, same namespace, same node).
	// It returns TRUE if the pod was inserted,
	// FALSE if it already existed or an error occurred
	Link(pod *core_v1.Pod) bool
	// Unlink removes the  pod from the list of monitored ones by this manager.
	// The second arguments specifies if the pod's firewall should be cleaned
	// or destroyed. It returns FALSE if the pod was not among the monitored ones,
	// and the number of remaining pods linked.
	Unlink(pod *core_v1.Pod, then UnlinkOperation) (bool, int)
	// LinkedPods returns a map of pods monitored by this firewall manager.
	LinkedPods() map[k8s_types.UID]string
	// IsPolicyEnforced returns true if this firewall enforces this policy
	IsPolicyEnforced(name string) bool
	// Name returns the name of this firewall manager
	Name() string
	// Selector returns the namespace and labels of the pods
	// monitored by this firewall manager
	Selector() (map[string]string, string)
	// EnforcePolicy enforces the provided policy.
	EnforcePolicy(policy pcn_types.ParsedPolicy, rules pcn_types.ParsedRules)
	// CeasePolicy will cease a policy, removing all rules generated by it
	// and won't react to pod events included by it anymore.
	CeasePolicy(policyName string)
	// Destroy destroys the current firewall manager.
	// This function should not be called manually, as it is called automatically
	// after a certain time has passed while monitoring no pods.
	// To destroy a particular firewall, see the Unlink function.
	Destroy()
}

// FirewallManager is the implementation of the firewall manager.
type FirewallManager struct {
	// ingressRules contains the ingress rules divided by policy
	ingressRules map[string][]k8sfirewall.ChainRule
	// egressRules contains the egress rules divided by policy
	egressRules map[string][]k8sfirewall.ChainRule
	// linkedPods is a map of pods monitored by this firewall manager
	linkedPods map[k8s_types.UID]string
	// Name is the name of this firewall manager
	name string
	// lock is firewall manager's main lock
	lock sync.Mutex
	// ingressDefaultAction is the default action for ingress
	ingressDefaultAction string
	// egressDefaultAction is the default action for egress
	egressDefaultAction string
	// ingressPoliciesCount is the count of ingress policies enforced
	ingressPoliciesCount int
	// egressPoliciesCount is the count of egress policies enforced
	egressPoliciesCount int
	// policyTypes is a map of policies types enforced.
	// Used to know how the default action should be handled.
	//policyTypes      map[string]string
	policyDirections map[string]string
	// policyActions contains a map of actions to be taken when a pod event occurs
	policyActions map[string]*subscriptions
	// selector defines what kind of pods this firewall is monitoring
	selector selector
	// priorities is the list of priorities
	priorities []policyPriority
}

// policyPriority is the priority of this policy.
// Most recently deployed policies take precedence over the older ones.
type policyPriority struct {
	policyName     string
	parentPriority int32
	priority       int32
	timestamp      time.Time
}

// selector is the selector for the pods this firewall manager is managing
type selector struct {
	namespace string
	labels    map[string]string
}

// subscriptions contains the templates (here called actions) that should be used
// when a pod event occurs, and the functions to be called when we need
// to unsubscribe.
type subscriptions struct {
	actions        map[string]*pcn_types.ParsedRules
	unsubscriptors []func()
}

// StartFirewall will start a new firewall manager
func StartFirewall(name, namespace string, labels map[string]string) PcnFirewallManager {
	logger.Infoln("Starting Firewall Manager, with name", name)

	manager := &FirewallManager{
		// Rules
		ingressRules: map[string][]k8sfirewall.ChainRule{},
		egressRules:  map[string][]k8sfirewall.ChainRule{},
		// The name (its key)
		name: name,
		// Selector
		selector: selector{
			namespace: namespace,
			labels:    labels,
		},
		// The counts
		ingressPoliciesCount: 0,
		egressPoliciesCount:  0,
		// Policy types
		//policyTypes: map[string]string{},
		policyDirections: map[string]string{},
		// Actions
		policyActions: map[string]*subscriptions{},
		// Linked pods
		linkedPods: map[k8s_types.UID]string{},
		// The default actions
		ingressDefaultAction: pcn_types.ActionDrop,
		egressDefaultAction:  pcn_types.ActionDrop,
		// The priorities
		priorities: []policyPriority{},
	}

	return manager
}

// Link adds a new pod to the list of pods that must be managed by this manager.
// Best practice is to only link similar pods
// (i.e.: same labels, same namespace, same node).
// It returns TRUE if the pod was inserted,
// FALSE if it already existed or an error occurred
func (d *FirewallManager) Link(pod *core_v1.Pod) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	f := "[FwManager-" + d.name + "](Link) "
	podIP := pod.Status.PodIP
	podUID := pod.UID
	name := "fw-" + podIP

	//-------------------------------------
	// Check firewall health and pod presence
	//-------------------------------------
	if ok, err := d.isFirewallOk(name); !ok {
		logger.Errorf(f+"Could not link firewall for pod %s: %s", name, err.Error())
		return false
	}
	_, alreadyLinked := d.linkedPods[podUID]
	if alreadyLinked {
		return false
	}

	//-------------------------------------
	// Extract the rules
	//-------------------------------------
	// We are going to get all rules regardless of the policy they belong to
	ingressRules := []k8sfirewall.ChainRule{}
	egressRules := []k8sfirewall.ChainRule{}

	if len(d.ingressRules) > 0 || len(d.egressRules) > 0 {
		// -- Ingress
		for _, rules := range d.ingressRules {
			ingressRules = append(ingressRules, rules...)
		}

		// -- Egress
		for _, rules := range d.egressRules {
			egressRules = append(egressRules, rules...)
		}
	}

	//-------------------------------------
	// Inject rules and change default actions
	//-------------------------------------
	if len(ingressRules) > 0 || len(egressRules) > 0 {
		if err := d.injecter(name, ingressRules, egressRules, nil, 0, 0); err != nil {
			// injecter fails only if pod's firewall is not ok (it is dying
			// or crashed or not found) so there's no point in going on.
			logger.Warningf(f+"Injecter encountered an error upon linking the pod: %s. Will stop here.", err)
			return false
		}
	}

	// -- ingress
	err := d.updateDefaultAction(name, "ingress", d.ingressDefaultAction)
	if err != nil {
		logger.Errorf(f+"Could not update the default ingress action: %s", err)
	} else {
		_, err := d.applyRules(name, "ingress")
		if err != nil {
			logger.Errorf(f+"Could not apply ingress rules: %s", err)
		}
	}

	// -- egress
	err = d.updateDefaultAction(name, "egress", d.egressDefaultAction)
	if err != nil {
		logger.Errorf(f+"Could not update the default egress action: %s", err)
	} else {
		_, err := d.applyRules(name, "egress")
		if err != nil {
			logger.Errorf(f+"Could not apply egress rules: %s", err)
		}
	}

	//-------------------------------------
	// Finally, link it
	//-------------------------------------
	// From now on, when this firewall manager will react to events,
	// this pod's firewall will be updated as well.
	d.linkedPods[podUID] = podIP
	logger.Infof(f+"Pod %s linked.", podIP)
	return true
}

// Unlink removes the  pod from the list of monitored ones by this manager.
// The second arguments specifies if the pod's firewall should be cleaned
// or destroyed. It returns FALSE if the pod was not among the monitored ones,
// and the number of remaining pods linked.
func (d *FirewallManager) Unlink(pod *core_v1.Pod, then UnlinkOperation) (bool, int) {
	d.lock.Lock()
	defer d.lock.Unlock()

	f := "[FwManager-" + d.name + "](Unlink) "
	podUID := pod.UID
	_, ok := d.linkedPods[podUID]
	if !ok {
		// This pod was not even linked
		return false, len(d.linkedPods)
	}

	podIP := d.linkedPods[pod.UID]
	name := "fw-" + podIP

	// Should I also destroy its firewall?
	switch then {
	case CleanFirewall:
		if i, e := d.cleanFw(name); i != nil || e != nil {
			logger.Warningln(f + "Could not properly clean firewall for the provided pod.")
		} else {
			d.updateDefaultAction(name, "ingress", pcn_types.ActionForward)
			d.applyRules(name, "ingress")
			d.updateDefaultAction(name, "egress", pcn_types.ActionForward)
			d.applyRules(name, "egress")
		}
	case DestroyFirewall:
		if err := d.destroyFw(name); err != nil {
			logger.Warningln(f+"Could not delete firewall for the provided pod:", err)
		}
	}

	delete(d.linkedPods, podUID)
	logger.Infof(f+"Pod %s unlinked.", podIP)
	return true, len(d.linkedPods)
}

// LinkedPods returns a map of pods monitored by this firewall manager.
func (d *FirewallManager) LinkedPods() map[k8s_types.UID]string {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.linkedPods
}

// Name returns the name of this firewall manager
func (d *FirewallManager) Name() string {
	return d.name
}

// EnforcePolicy enforces the provided policy.
func (d *FirewallManager) EnforcePolicy(policy pcn_types.ParsedPolicy, rules pcn_types.ParsedRules) {
	// Only one policy at a time, please
	d.lock.Lock()
	defer d.lock.Unlock()
	f := "[FwManager-" + d.name + "](EnforcePolicy) "

	if _, exists := d.policyDirections[policy.Name]; exists {
		logger.Errorf(f+"Policy %s is already enforced. Will stop here.", policy.Name)
		return
	}
	logger.Infof(f+"Enforcing policy %s...", policy.Name)

	//-------------------------------------
	// Define the actions
	//-------------------------------------

	d.definePolicyActions(policy)

	//-------------------------------------
	// Store the rules
	//-------------------------------------

	ingressRules, egressRules := d.storeRules(policy.Name, "", rules.Outgoing, rules.Incoming)

	//-------------------------------------
	// Update default actions
	//-------------------------------------

	// update the policy type, so that later - if this policy is removed -
	// we can enforce isolation mode correctly
	d.policyDirections[policy.Name] = policy.Direction
	d.updateCounts("increase", dirToChain(policy.Direction))

	//-------------------------------------
	// Set its priority
	//-------------------------------------

	priority := calculatePolicyOffset(policy, d.priorities)
	d.insertPriority(policy, priority)

	// By setting its priority, we know where to start injecting rules from
	iStartFrom, eStartFrom := d.calculateInsertionIDs(policy.Name)

	//-------------------------------------
	// Inject the rules on each firewall
	//-------------------------------------

	if len(d.linkedPods) == 0 {
		logger.Infoln(f + "There are no linked pods. Stopping here.")
		return
	}

	var injectWaiter sync.WaitGroup
	injectWaiter.Add(len(d.linkedPods))

	for _, ip := range d.linkedPods {
		name := "fw-" + ip
		go d.injecter(name, ingressRules, egressRules, &injectWaiter, iStartFrom, eStartFrom)
	}
	injectWaiter.Wait()
}

// insertPriority insertes this policy in the priorities list
func (d *FirewallManager) insertPriority(policy pcn_types.ParsedPolicy, offset int) {
	newPriority := []policyPriority{
		policyPriority{
			policyName:     policy.Name,
			parentPriority: policy.ParentPolicy.Priority,
			priority:       policy.Priority,
			timestamp:      policy.CreationTime.Time,
		},
	}

	//https://github.com/golang/go/wiki/SliceTrick
	d.priorities = append(d.priorities[:offset], append(newPriority, d.priorities[offset:]...)...)
}

// popPriority removes a policy from the list of the priorities
func (d *FirewallManager) popPriority(policyName string) {
	offset := -1

	// Find it
	for i, p := range d.priorities {
		if p.policyName == policyName {
			offset = i
		}
	}

	if offset == -1 {
		logger.Errorf("[FwManager-"+d.name+"])(popPriority) Policy with name %s was not found in priorities list. Will stop here.", policyName)
		return
	}

	d.priorities = append(d.priorities[:offset], d.priorities[offset+1:]...)
}

// calculateInsertionIDs gets the first useful rules IDs for insertion,
// based on the policy's priority
func (d *FirewallManager) calculateInsertionIDs(policyName string) (int32, int32) {
	iStartFrom := 0
	eStartFrom := 0

	/* Example: we want to insert 2 rules for policy my-policy, this is the current list of priorities:
	0) my-other-policy (it has 2 rules)
	2) my-policy (it has 1 rule)
	3) yet-another-policy (it has 2 rules)

	On each iteration we count the number of rules for that policy, and we stop when the policy is the one we need.
	First iteration: 2 rules -> 0 + 2 = start from id 2 -> the policy is not the one we need -> go on.
	Second iteration: 1 rule -> 2 + 1 = start from id 3 -> the policy is ok > stop.
	We're going to insert from id 3.
	*/

	for _, currentPolicy := range d.priorities {
		if currentPolicy.policyName == policyName {
			break
		}

		// jump the rules
		iStartFrom += len(d.ingressRules[currentPolicy.policyName])
		eStartFrom += len(d.egressRules[currentPolicy.policyName])
	}

	return int32(iStartFrom), int32(eStartFrom)
}

// updateCounts updates the internal counts of policies types enforced,
// making sure default actions are respected. This is just a convenient method
// used to keep core methods (EnforcePolicy and CeasePolicy) as clean and
// as readable as possible. When possible, this function is used in place of
// increaseCount or decreaseCount, as it is preferrable to do it like this.
func (d *FirewallManager) updateCounts(operation, policyType string) {
	// BRIEF: read increaseCounts and decreaseCounts for an explanation
	// of when and why these functions are called.

	f := "[FwManager-" + d.name + "](updateCounts) "

	//-------------------------------------
	// Increase
	//-------------------------------------

	increase := func() {
		directions := []string{}

		// -- Increase the counts and append the directions to update accordingly.
		if (policyType == "ingress" || policyType == "*") && d.increaseCount("ingress") {
			directions = append(directions, "ingress")
		}
		if (policyType == "egress" || policyType == "*") && d.increaseCount("egress") {
			directions = append(directions, "egress")
		}

		if len(directions) == 0 {
			return
		}

		// -- Let's now update the default actions.
		for _, ip := range d.linkedPods {
			name := "fw-" + ip
			for _, direction := range directions {
				err := d.updateDefaultAction(name, direction, pcn_types.ActionDrop)
				if err != nil {
					logger.Errorf(f+"Could not update default action for firewall %s: %s", name, direction)
				} else {
					if _, err := d.applyRules(name, direction); err != nil {
						logger.Errorf(f+"Could not apply rules for firewall %s: %s", name, direction)
					}
				}
			}
		}
	}

	//-------------------------------------
	// Decrease
	//-------------------------------------

	decrease := func() {
		directions := []string{}

		// -- Decrease the counts and append the directions to update accordingly.
		if (policyType == "ingress" || policyType == "*") && d.decreaseCount("ingress") {
			directions = append(directions, "ingress")
		}
		if (policyType == "egress" || policyType == "*") && d.decreaseCount("egress") {
			directions = append(directions, "egress")
		}

		if len(directions) == 0 {
			return
		}

		// -- Let's now update the default actions.
		for _, ip := range d.linkedPods {
			name := "fw-" + ip
			for _, direction := range directions {
				err := d.updateDefaultAction(name, direction, pcn_types.ActionForward)
				if err != nil {
					logger.Errorf(f+"Could not update default action for firewall %s: %s", name, direction)
				} else {
					if _, err := d.applyRules(name, direction); err != nil {
						logger.Errorf(f+"Could not apply rules for firewall %s: %s", name, direction)
					}
				}
			}
		}
	}

	switch operation {
	case "increase":
		increase()
	case "decrease":
		decrease()
	}
}

// increaseCount increases the count of policies enforced
// and changes the default action for the provided direction, if needed.
// It returns TRUE if the corresponding action should be updated
func (d *FirewallManager) increaseCount(which string) bool {
	// BRIEF: this function is called when a new policy is deployed
	// with the appropriate direction.
	// If there are no policies, the default action is FORWARD.
	// If there is at least one, then the default action should be updated to DROP
	// because only what is allowed is forwarded. This function returns true when
	// there is only one policy, because that's when we should actually switch
	// to DROP (we were in FORWARD)

	// Ingress
	if which == "ingress" {
		d.ingressPoliciesCount++

		if d.ingressPoliciesCount > 0 {
			d.ingressDefaultAction = pcn_types.ActionDrop
			// If this is the *first* ingress policy, then switch to drop,
			// otherwise no need to do that (it's already DROP)
			if d.ingressPoliciesCount == 1 {
				return true
			}
		}
	}

	// Egress
	if which == "egress" {
		d.egressPoliciesCount++

		if d.egressPoliciesCount > 0 {
			d.egressDefaultAction = pcn_types.ActionDrop

			if d.egressPoliciesCount == 1 {
				return true
			}
		}
	}

	return false
}

// decreaseCount decreases the count of policies enforced and changes
// the default action for the provided direction, if needed.
// It returns TRUE if the corresponding action should be updated
func (d *FirewallManager) decreaseCount(which string) bool {
	f := "[FwManager-" + d.name + "](decreaseCount) "
	// BRIEF: this function is called when a policy must be ceased.
	// If - after ceasing it - we have no policies enforced,
	// then the default action must be FORWARD.
	// If there is at least one, then the default action should remain DROP
	// This function returns true when there are no policies enforced,
	// because that's when we should actually switch to FORWARD (we were in DROP)

	if which == "ingress" {
		if d.ingressPoliciesCount == 0 {
			logger.Errorf(f + "Cannot further decrease ingress policies count (it is 0). Going to stop here.")
			return true
		}

		d.ingressPoliciesCount--
		// Return to default=FORWARD only if there are no policies anymore
		// after removing this
		if d.ingressPoliciesCount == 0 {
			d.ingressDefaultAction = pcn_types.ActionForward
			return true
		}
	}

	if which == "egress" {
		if d.egressPoliciesCount == 0 {
			logger.Errorf(f + "Cannot further decrease egress policies count (it is 0). Going to stop here.")
			return true
		}
		if d.egressPoliciesCount--; d.egressPoliciesCount == 0 {
			d.egressDefaultAction = pcn_types.ActionForward
			return true
		}
	}

	return false
}

// storeRules stores rules in memory according to their policy
func (d *FirewallManager) storeRules(policyName, target string, ingress, egress []k8sfirewall.ChainRule) ([]k8sfirewall.ChainRule, []k8sfirewall.ChainRule) {
	if _, exists := d.ingressRules[policyName]; !exists {
		d.ingressRules[policyName] = []k8sfirewall.ChainRule{}
	}
	if _, exists := d.egressRules[policyName]; !exists {
		d.egressRules[policyName] = []k8sfirewall.ChainRule{}
	}

	description := "policy=" + policyName
	newIngress := make([]k8sfirewall.ChainRule, len(ingress))
	newEgress := make([]k8sfirewall.ChainRule, len(egress))

	// --- ingress
	for i, rule := range ingress {
		newIngress[i] = rule
		newIngress[i].Description = description
		if len(target) > 0 {
			newIngress[i].Dst = target
		}

		d.ingressRules[policyName] = append(d.ingressRules[policyName], newIngress[i])
	}

	// -- Egress
	for i, rule := range egress {
		newEgress[i] = rule
		newEgress[i].Description = description
		if len(target) > 0 {
			newEgress[i].Src = target
		}

		d.egressRules[policyName] = append(d.egressRules[policyName], newEgress[i])
	}
	return newIngress, newEgress
}

// injecter is a convenient method for injecting rules for a single firewall
// for both directions
func (d *FirewallManager) injecter(firewall string, ingressRules, egressRules []k8sfirewall.ChainRule, waiter *sync.WaitGroup, iStartFrom, eStartFrom int32) error {
	f := "[FwManager-" + d.name + "](injecter) "
	// Should I notify caller when I'm done?
	if waiter != nil {
		defer waiter.Done()
	}

	// Is firewall ok?
	if ok, err := d.isFirewallOk(firewall); !ok {
		logger.Errorf(f+"Could not inject rules. Firewall is not ok: %s", err)
		return err
	}

	//-------------------------------------
	// Inject rules direction concurrently
	//-------------------------------------
	var injectWaiter sync.WaitGroup
	injectWaiter.Add(2)
	defer injectWaiter.Wait()

	go d.injectRules(firewall, "ingress", ingressRules, &injectWaiter, iStartFrom)
	go d.injectRules(firewall, "egress", egressRules, &injectWaiter, eStartFrom)

	return nil
}

// injectRules is a wrapper for firewall's CreateFirewallChainRuleListByID
// and CreateFirewallChainApplyRulesByID methods.
func (d *FirewallManager) injectRules(firewall, direction string, rules []k8sfirewall.ChainRule, waiter *sync.WaitGroup, startFrom int32) error {
	f := "[FwManager-" + d.name + "](injectRules) "
	// Should I notify caller when I'm done?
	if waiter != nil {
		defer waiter.Done()
	}

	//-------------------------------------
	// Inject & apply
	//-------------------------------------
	// The ip of the pod we are protecting. Used for the SRC or the DST
	me := strings.Split(firewall, "-")[1]

	// We are using the insert call here, which adds the rule on the startFrom id
	// and pushes the other rules downwards. In order to preserve original order,
	// we're going to start injecting from the last to the first.

	len := len(rules)
	for i := len - 1; i > -1; i-- {
		ruleToInsert := k8sfirewall.ChainInsertInput(rules[i])
		ruleToInsert.Id = startFrom

		// There is only the pod on the other side of the link, so all packets
		// travelling there are obviously either going to or from the pod.
		// So there is no need to include the pod as a destination or source
		// in the rules. But it helps to keep them clearer and precise.
		if direction == "ingress" {
			ruleToInsert.Src = me
		} else {
			ruleToInsert.Dst = me
		}

		_, response, err := fwAPI.CreateFirewallChainInsertByID(nil, firewall, direction, ruleToInsert)
		if err != nil {
			logger.Errorf(f+"Error while trying to inject rule: %s, %+v", err, response)
			// This rule had an error, but we still gotta push the other ones dude...
			//return err
		}
	}

	// Now apply the changes
	if response, err := d.applyRules(firewall, direction); err != nil {
		logger.Errorf(f+"Error while trying to apply rules: %s, %+v", err, response)
		return err
	}

	return nil
}

// definePolicyActions subscribes to the appropriate events
// and defines the actions to be taken when that event happens.
func (d *FirewallManager) definePolicyActions(policy pcn_types.ParsedPolicy) {
	f := "[FwManager-" + d.name + "](definePolicyActions) "
	if len(policy.Peer.Key) == 0 {
		logger.Infoln(f + "Policy does not need to react to events. Stopping here.")
		return
	}

	shouldSubscribe := false
	if _, exists := d.policyActions[policy.Peer.Key]; !exists {
		d.policyActions[policy.Peer.Key] = &subscriptions{
			actions: map[string]*pcn_types.ParsedRules{},
		}
		shouldSubscribe = true
	}

	// Define the action...
	if _, exists := d.policyActions[policy.Peer.Key].actions[policy.Name]; !exists {
		d.policyActions[policy.Peer.Key].actions[policy.Name] = &pcn_types.ParsedRules{}
	}

	d.policyActions[policy.Peer.Key].actions[policy.Name].Incoming = append(d.policyActions[policy.Peer.Key].actions[policy.Name].Incoming, policy.Templates.Incoming...)
	d.policyActions[policy.Peer.Key].actions[policy.Name].Outgoing = append(d.policyActions[policy.Peer.Key].actions[policy.Name].Outgoing, policy.Templates.Outgoing...)

	if shouldSubscribe {
		// Prepare the subscription queries
		podQuery := policy.Peer.Peer
		nsQuery := policy.Peer.Namespace

		// Finally, susbcribe...
		// NOTE: this function will handle both deletion and updates. This
		// is because when a pod is deleted, it doesn't have its IP set
		// anymore, so we cannot subscribe to DELETE events but rather to
		// UPDATEs and manually check if it is terminating.
		unsub, err := pcn_controllers.Pods().Subscribe(pcn_types.Update, podQuery, nsQuery, nil, pcn_types.PodAnyPhase, func(pod *core_v1.Pod, prev *core_v1.Pod) {
			// Phase is not running?
			if pod.Status.Phase != pcn_types.PodRunning {
				return
			}

			// What to do?
			event := pcn_types.Update

			// Is the pod terminating?
			if pod.ObjectMeta.DeletionTimestamp != nil {
				event = pcn_types.Delete
			}

			d.reactToPod(pcn_types.EventType(event), pod, prev, policy.Peer.Key, podQuery)
		})

		if err == nil {
			d.policyActions[policy.Peer.Key].unsubscriptors = append(d.policyActions[policy.Peer.Key].unsubscriptors, unsub)
		} else {
			logger.Errorf(f+"Could not subscribe to changes! %s", err)
		}
	}
}

// reactToPod is called whenever a monitored pod event occurs.
// E.g.: I should accept connections from Pod A, and a new Pod A is born.
// This function knows what to do when that event happens.
func (d *FirewallManager) reactToPod(event pcn_types.EventType, pod *core_v1.Pod, prev *core_v1.Pod, actionKey string, podQuery *pcn_types.ObjectQuery) {
	d.lock.Lock()
	defer d.lock.Unlock()

	f := "[FwManager-" + d.name + "](reactToPod) "

	if len(pod.Status.PodIP) == 0 {
		return
	}

	//-------------------------------------
	//	Delete
	//-------------------------------------

	del := func(ip string) {
		virtualIP := utils.GetPodVirtualIP(ip)
		logger.Infof("Reacting to deleted pod: name %s, IP %s, labels %v, namespace %s", pod.Name, pod.Status.PodIP, pod.Labels, pod.Namespace)

		deleteIngress := func() {
			//	For each policy, go get all the rules in which this ip was present.
			rulesToDelete := []k8sfirewall.ChainRule{}
			rulesToKeep := []k8sfirewall.ChainRule{}
			for policy, rules := range d.ingressRules {
				for _, rule := range rules {
					if rule.Dst == ip || rule.Dst == virtualIP {
						rulesToDelete = append(rulesToDelete, rule)
					} else {
						rulesToKeep = append(rulesToKeep, rule)
					}
				}

				d.ingressRules[policy] = rulesToKeep
			}

			if len(rulesToDelete) == 0 {
				logger.Debugln(f + "No rules to delete in ingress")
				return
			}

			//	Delete the rules on each linked pod
			for _, fwIP := range d.linkedPods {
				name := "fw-" + fwIP
				d.deleteRules(name, "ingress", rulesToDelete)
				d.applyRules(name, "ingress")
			}
		}

		//	--- Egress
		deleteEgress := func() {
			rulesToDelete := []k8sfirewall.ChainRule{}
			rulesToKeep := []k8sfirewall.ChainRule{}
			for policy, rules := range d.egressRules {
				for _, rule := range rules {
					if rule.Src == ip || rule.Src == virtualIP {
						rulesToDelete = append(rulesToDelete, rule)
					} else {
						rulesToKeep = append(rulesToKeep, rule)
					}
				}

				d.egressRules[policy] = rulesToKeep
			}

			if len(rulesToDelete) == 0 {
				logger.Debugln(f + "No rules to delete in egress")
				return
			}

			for _, fwIP := range d.linkedPods {
				name := "fw-" + fwIP
				d.deleteRules(name, "egress", rulesToDelete)
				d.applyRules(name, "egress")
			}
		}

		deleteIngress()
		deleteEgress()
	}

	//-------------------------------------
	//	Update
	//-------------------------------------

	upd := func(ip string) {
		logger.Infof(f+"Reacting to updated pod: name %s, IP %s, labels %v, namespace %s", pod.Name, pod.Status.PodIP, pod.Labels, pod.Namespace)

		//	Basic checks
		actions, exist := d.policyActions[actionKey]
		if !exist {
			logger.Warningln(f + "Could not find any actions with this key")
			return
		}
		if len(actions.actions) == 0 {
			logger.Warningln(f + "There are no actions to be taken!")
			return
		}

		// Check for any changes. If the pod changed the labels, then we have
		// to actually remove this, not add it
		if podQuery != nil && len(pod.Labels) > 0 {
			if !utils.AreLabelsContained(podQuery.Labels, pod.Labels) {
				del(pod.Status.PodIP)
				return
			}
		}

		//	Build rules according to the policy's priority
		for policy, rules := range actions.actions {

			ingress := []k8sfirewall.ChainRule{}
			egress := []k8sfirewall.ChainRule{}

			// first calculate the priority
			iStartFrom, eStartFrom := d.calculateInsertionIDs(policy)
			podIPs := []string{ip, utils.GetPodVirtualIP(ip)}

			//	Then format the rules
			for _, podIP := range podIPs {
				ingressRules, egressRules := d.storeRules(policy, podIP, rules.Outgoing, rules.Incoming)
				ingress = append(ingress, ingressRules...)
				egress = append(egress, egressRules...)
			}

			//	Now inject the rules in all firewalls linked.
			//	This usually is a matter of 1-2 rules,
			// so no need to do this in a separate goroutine.
			for _, f := range d.linkedPods {
				name := "fw-" + f
				d.injecter(name, ingress, egress, nil, iStartFrom, eStartFrom)
			}
		}
	}

	//-------------------------------------
	//	Main entrypoint: what to do?
	//-------------------------------------

	switch event {
	case pcn_types.Update:
		upd(pod.Status.PodIP)
	case pcn_types.Delete:
		del(pod.Status.PodIP)
	}
}

// deleteAllPolicyRules deletes all rules mentioned in a policy
func (d *FirewallManager) deleteAllPolicyRules(policy string) {
	//-------------------------------------
	// Ingress
	//-------------------------------------

	func() {
		defer delete(d.ingressRules, policy)

		rules := d.ingressRules[policy]
		if len(rules) == 0 {
			return
		}

		// Delete the found rules from each linked pod
		for _, ip := range d.linkedPods {
			name := "fw-" + ip
			d.deleteRules(name, "ingress", rules)
			d.applyRules(name, "ingress")
		}
	}()

	//-------------------------------------
	// Egress
	//-------------------------------------

	func() {
		defer delete(d.egressRules, policy)
		//defer waiter.Done()

		rules := d.egressRules[policy]
		if len(rules) == 0 {
			return
		}

		// Delete the found rules from each linked pod
		for _, ip := range d.linkedPods {
			name := "fw-" + ip
			d.deleteRules(name, "egress", rules)
			d.applyRules(name, "egress")
		}
	}()
}

// deletePolicyActions delete all templates generated by a specific policy.
// So that the firewall manager will not generate those rules anymore
// when it will react to a certain pod.
func (d *FirewallManager) deletePolicyActions(policy string) {
	flaggedForDeletion := []string{}

	//-------------------------------------
	// Delete this policy from the actions
	//-------------------------------------
	for key, action := range d.policyActions {
		delete(action.actions, policy)

		// This action belongs to no policies anymore?
		if len(action.actions) == 0 {
			flaggedForDeletion = append(flaggedForDeletion, key)
		}
	}

	//-------------------------------------
	// Delete actions with no policies
	//-------------------------------------
	// If, after deleting the policy from the actions, the action has no
	// policies anymore then we need to stop monitoring that pod!
	for _, flaggedKey := range flaggedForDeletion {
		for _, unsubscribe := range d.policyActions[flaggedKey].unsubscriptors {
			unsubscribe()
		}

		delete(d.policyActions, flaggedKey)
	}
}

// CeasePolicy will cease a policy, removing all rules generated by it
// and won't react to pod events included by it anymore.
func (d *FirewallManager) CeasePolicy(policyName string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	f := "[FwManager-" + d.name + "](CeasePolicy) "

	if _, exists := d.policyDirections[policyName]; !exists {
		logger.Errorf(f+"Policy %s is not currently enforced. Will stop here.", policyName)
		return
	}

	//-------------------------------------
	// Delete all rules generated by this policy
	//-------------------------------------

	d.deleteAllPolicyRules(policyName)

	//-------------------------------------
	// Remove this policy's templates from the actions
	//-------------------------------------

	d.deletePolicyActions(policyName)

	//-------------------------------------
	// Remove this policy's priority
	//-------------------------------------

	d.popPriority(policyName)

	//-------------------------------------
	// Update the default actions
	//-------------------------------------
	// So we just ceased a policy, we now need to update the default actions
	if _, exists := d.policyDirections[policyName]; exists {
		policyDir := d.policyDirections[policyName]
		d.updateCounts("decrease", dirToChain(policyDir))
		delete(d.policyDirections, policyName)
	} else {
		logger.Warningln(f+policyName, "was not listed among policy types!")
	}
}

// deleteRules is a wrapper for DeleteFirewallChainRuleByID method,
// deleting multiple rules.
func (d *FirewallManager) deleteRules(fw, direction string, rules []k8sfirewall.ChainRule) error {
	me := strings.Split(fw, "-")[1]
	f := "[FwManager-" + d.name + "](deleteRules) "

	// this is a fake deep copy-cast.
	cast := func(rule k8sfirewall.ChainRule) k8sfirewall.ChainDeleteInput {

		src := rule.Src
		dst := rule.Dst

		if direction == "ingress" {
			src = me
		} else {
			dst = me
		}

		return k8sfirewall.ChainDeleteInput{
			Src:         src,
			Dst:         dst,
			L4proto:     rule.L4proto,
			Sport:       rule.Sport,
			Dport:       rule.Dport,
			Tcpflags:    rule.Tcpflags,
			Conntrack:   rule.Conntrack,
			Action:      rule.Action,
			Description: rule.Description,
		}
	}

	// No need to do this with separate threads...
	for _, rule := range rules {
		// Delete the rule not by its ID, but by the fields it is composed of.
		response, err := fwAPI.CreateFirewallChainDeleteByID(nil, fw, direction, cast(rule))
		if err != nil {
			logger.Errorf(f+"Error while trying to delete this rule: %+v, in %s for firewall %s. Error %s, response: %+v", rule, direction, fw, err.Error(), response)
		}
	}

	return nil
}

// IsPolicyEnforced returns true if this firewall enforces this policy
func (d *FirewallManager) IsPolicyEnforced(name string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, exists := d.policyDirections[name]
	return exists
}

// Selector returns the namespace and labels of the pods
// monitored by this firewall manager
func (d *FirewallManager) Selector() (map[string]string, string) {
	return d.selector.labels, d.selector.namespace
}

// isFirewallOk checks if the firewall is ok.
// Used to check if firewall exists and is healthy.
func (d *FirewallManager) isFirewallOk(firewall string) (bool, error) {
	// We are going to do that by reading its uuid
	if _, _, err := fwAPI.ReadFirewallUuidByID(nil, firewall); err != nil {
		return false, err
	}
	return true, nil
}

// updateDefaultAction is a wrapper for UpdateFirewallChainDefaultByID method.
func (d *FirewallManager) updateDefaultAction(firewall, direction, to string) error {
	// To is enclosed between \" because otherwise the swagger-generated API
	// would wrongly send forward instead of "forward".
	actualTo := "\"" + to + "\""
	_, err := fwAPI.UpdateFirewallChainDefaultByID(nil, firewall, direction, actualTo)
	return err
}

// destroyFw destroy a firewall linked by this firewall manager
func (d *FirewallManager) destroyFw(name string) error {
	_, err := fwAPI.DeleteFirewallByID(nil, name)
	return err
}

// cleanFw cleans the firewall linked by this firewall manager
func (d *FirewallManager) cleanFw(name string) (error, error) {
	var iErr error
	var eErr error

	if _, err := fwAPI.DeleteFirewallChainRuleListByID(nil, name, "ingress"); err != nil {
		iErr = err
	}
	if _, err := fwAPI.DeleteFirewallChainRuleListByID(nil, name, "egress"); err != nil {
		eErr = err
	}

	return iErr, eErr
}

// applyRules is a wrapper for CreateFirewallChainApplyRulesByID method.
func (d *FirewallManager) applyRules(firewall, direction string) (bool, error) {
	out, _, err := fwAPI.CreateFirewallChainApplyRulesByID(nil, firewall, direction)
	return out.Result, err
}

// Destroy destroys the current firewall manager.
// This function should not be called manually, as it is called automatically
// after a certain time has passed while monitoring no pods.
// To destroy a particular firewall, see the Unlink function.
func (d *FirewallManager) Destroy() {
	d.lock.Lock()
	defer d.lock.Unlock()

	//-------------------------------------
	//	Unsubscribe from all actions
	//-------------------------------------
	keysToDelete := make([]string, len(d.policyActions))
	i := 0

	//	-- Unsubscribe
	for key, action := range d.policyActions {
		for _, unsubscribe := range action.unsubscriptors {
			unsubscribe()
		}
		keysToDelete[i] = key
		i++
	}

	//	-- Delete the action.
	//	We do this so that queued actions will instantly return with no harm.
	for _, key := range keysToDelete {
		delete(d.policyActions, key)
	}

	logger.Infoln("[FwManager-" + d.name + "] Good bye!")
}
