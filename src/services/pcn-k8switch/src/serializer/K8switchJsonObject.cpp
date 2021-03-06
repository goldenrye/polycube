/**
* k8switch API
* k8switch API generated from k8switch.yang
*
* OpenAPI spec version: 1.0.0
*
* NOTE: This class is auto generated by the swagger code generator program.
* https://github.com/polycube-network/swagger-codegen.git
* branch polycube
*/


/* Do not edit this file manually */



#include "K8switchJsonObject.h"
#include <regex>

namespace io {
namespace swagger {
namespace server {
namespace model {

K8switchJsonObject::K8switchJsonObject() {
  m_nameIsSet = false;
  m_portsIsSet = false;
  m_clusterIpSubnetIsSet = false;
  m_clientSubnetIsSet = false;
  m_virtualClientSubnetIsSet = false;
  m_serviceIsSet = false;
  m_fwdTableIsSet = false;
}

K8switchJsonObject::K8switchJsonObject(const nlohmann::json &val) :
  JsonObjectBase(val) {
  m_nameIsSet = false;
  m_portsIsSet = false;
  m_clusterIpSubnetIsSet = false;
  m_clientSubnetIsSet = false;
  m_virtualClientSubnetIsSet = false;
  m_serviceIsSet = false;
  m_fwdTableIsSet = false;


  if (val.count("name")) {
    setName(val.at("name").get<std::string>());
  }

  if (val.count("ports")) {
    for (auto& item : val["ports"]) {
      PortsJsonObject newItem{ item };
      m_ports.push_back(newItem);
    }

    m_portsIsSet = true;
  }

  if (val.count("cluster-ip-subnet")) {
    setClusterIpSubnet(val.at("cluster-ip-subnet").get<std::string>());
  }

  if (val.count("client-subnet")) {
    setClientSubnet(val.at("client-subnet").get<std::string>());
  }

  if (val.count("virtual-client-subnet")) {
    setVirtualClientSubnet(val.at("virtual-client-subnet").get<std::string>());
  }

  if (val.count("service")) {
    for (auto& item : val["service"]) {
      ServiceJsonObject newItem{ item };
      m_service.push_back(newItem);
    }

    m_serviceIsSet = true;
  }

  if (val.count("fwd-table")) {
    for (auto& item : val["fwd-table"]) {
      FwdTableJsonObject newItem{ item };
      m_fwdTable.push_back(newItem);
    }

    m_fwdTableIsSet = true;
  }
}

nlohmann::json K8switchJsonObject::toJson() const {
  nlohmann::json val = nlohmann::json::object();
  if (!getBase().is_null()) {
    val.update(getBase());
  }

  if (m_nameIsSet) {
    val["name"] = m_name;
  }

  {
    nlohmann::json jsonArray;
    for (auto& item : m_ports) {
      jsonArray.push_back(JsonObjectBase::toJson(item));
    }

    if (jsonArray.size() > 0) {
      val["ports"] = jsonArray;
    }
  }

  if (m_clusterIpSubnetIsSet) {
    val["cluster-ip-subnet"] = m_clusterIpSubnet;
  }

  if (m_clientSubnetIsSet) {
    val["client-subnet"] = m_clientSubnet;
  }

  if (m_virtualClientSubnetIsSet) {
    val["virtual-client-subnet"] = m_virtualClientSubnet;
  }

  {
    nlohmann::json jsonArray;
    for (auto& item : m_service) {
      jsonArray.push_back(JsonObjectBase::toJson(item));
    }

    if (jsonArray.size() > 0) {
      val["service"] = jsonArray;
    }
  }

  {
    nlohmann::json jsonArray;
    for (auto& item : m_fwdTable) {
      jsonArray.push_back(JsonObjectBase::toJson(item));
    }

    if (jsonArray.size() > 0) {
      val["fwd-table"] = jsonArray;
    }
  }

  return val;
}

std::string K8switchJsonObject::getName() const {
  return m_name;
}

void K8switchJsonObject::setName(std::string value) {
  m_name = value;
  m_nameIsSet = true;
}

bool K8switchJsonObject::nameIsSet() const {
  return m_nameIsSet;
}



const std::vector<PortsJsonObject>& K8switchJsonObject::getPorts() const{
  return m_ports;
}

void K8switchJsonObject::addPorts(PortsJsonObject value) {
  m_ports.push_back(value);
  m_portsIsSet = true;
}


bool K8switchJsonObject::portsIsSet() const {
  return m_portsIsSet;
}

void K8switchJsonObject::unsetPorts() {
  m_portsIsSet = false;
}

std::string K8switchJsonObject::getClusterIpSubnet() const {
  return m_clusterIpSubnet;
}

void K8switchJsonObject::setClusterIpSubnet(std::string value) {
  m_clusterIpSubnet = value;
  m_clusterIpSubnetIsSet = true;
}

bool K8switchJsonObject::clusterIpSubnetIsSet() const {
  return m_clusterIpSubnetIsSet;
}



std::string K8switchJsonObject::getClientSubnet() const {
  return m_clientSubnet;
}

void K8switchJsonObject::setClientSubnet(std::string value) {
  m_clientSubnet = value;
  m_clientSubnetIsSet = true;
}

bool K8switchJsonObject::clientSubnetIsSet() const {
  return m_clientSubnetIsSet;
}



std::string K8switchJsonObject::getVirtualClientSubnet() const {
  return m_virtualClientSubnet;
}

void K8switchJsonObject::setVirtualClientSubnet(std::string value) {
  m_virtualClientSubnet = value;
  m_virtualClientSubnetIsSet = true;
}

bool K8switchJsonObject::virtualClientSubnetIsSet() const {
  return m_virtualClientSubnetIsSet;
}



const std::vector<ServiceJsonObject>& K8switchJsonObject::getService() const{
  return m_service;
}

void K8switchJsonObject::addService(ServiceJsonObject value) {
  m_service.push_back(value);
  m_serviceIsSet = true;
}


bool K8switchJsonObject::serviceIsSet() const {
  return m_serviceIsSet;
}

void K8switchJsonObject::unsetService() {
  m_serviceIsSet = false;
}

const std::vector<FwdTableJsonObject>& K8switchJsonObject::getFwdTable() const{
  return m_fwdTable;
}

void K8switchJsonObject::addFwdTable(FwdTableJsonObject value) {
  m_fwdTable.push_back(value);
  m_fwdTableIsSet = true;
}


bool K8switchJsonObject::fwdTableIsSet() const {
  return m_fwdTableIsSet;
}

void K8switchJsonObject::unsetFwdTable() {
  m_fwdTableIsSet = false;
}


}
}
}
}

