# \SlbApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSlbByID**](SlbApi.md#CreateSlbByID) | **Post** /slb/{name}/ | Create slb by ID
[**DeleteSlbByID**](SlbApi.md#DeleteSlbByID) | **Delete** /slb/{name}/ | Delete slb by ID
[**ReadSlbByID**](SlbApi.md#ReadSlbByID) | **Get** /slb/{name}/ | Read slb by ID
[**ReadSlbChannelLenByID**](SlbApi.md#ReadSlbChannelLenByID) | **Get** /slb/{name}/channel-len/ | Read channel-len by ID
[**ReadSlbChannelLocByID**](SlbApi.md#ReadSlbChannelLocByID) | **Get** /slb/{name}/channel-loc/ | Read channel-loc by ID
[**ReadSlbEgressActionByID**](SlbApi.md#ReadSlbEgressActionByID) | **Get** /slb/{name}/egress-action/ | Read egress-action by ID
[**ReadSlbIngressActionByID**](SlbApi.md#ReadSlbIngressActionByID) | **Get** /slb/{name}/ingress-action/ | Read ingress-action by ID
[**ReadSlbListByID**](SlbApi.md#ReadSlbListByID) | **Get** /slb/ | Read slb by ID
[**ReadSlbLoglevelByID**](SlbApi.md#ReadSlbLoglevelByID) | **Get** /slb/{name}/loglevel/ | Read loglevel by ID
[**ReadSlbServerIdByID**](SlbApi.md#ReadSlbServerIdByID) | **Get** /slb/{name}/server-id/ | Read server-id by ID
[**ReadSlbServiceNameByID**](SlbApi.md#ReadSlbServiceNameByID) | **Get** /slb/{name}/service-name/ | Read service-name by ID
[**ReadSlbTypeByID**](SlbApi.md#ReadSlbTypeByID) | **Get** /slb/{name}/type/ | Read type by ID
[**ReadSlbUuidByID**](SlbApi.md#ReadSlbUuidByID) | **Get** /slb/{name}/uuid/ | Read uuid by ID
[**ReplaceSlbByID**](SlbApi.md#ReplaceSlbByID) | **Put** /slb/{name}/ | Replace slb by ID
[**UpdateSlbByID**](SlbApi.md#UpdateSlbByID) | **Patch** /slb/{name}/ | Update slb by ID
[**UpdateSlbChannelLenByID**](SlbApi.md#UpdateSlbChannelLenByID) | **Patch** /slb/{name}/channel-len/ | Update channel-len by ID
[**UpdateSlbChannelLocByID**](SlbApi.md#UpdateSlbChannelLocByID) | **Patch** /slb/{name}/channel-loc/ | Update channel-loc by ID
[**UpdateSlbEgressActionByID**](SlbApi.md#UpdateSlbEgressActionByID) | **Patch** /slb/{name}/egress-action/ | Update egress-action by ID
[**UpdateSlbIngressActionByID**](SlbApi.md#UpdateSlbIngressActionByID) | **Patch** /slb/{name}/ingress-action/ | Update ingress-action by ID
[**UpdateSlbListByID**](SlbApi.md#UpdateSlbListByID) | **Patch** /slb/ | Update slb by ID
[**UpdateSlbLoglevelByID**](SlbApi.md#UpdateSlbLoglevelByID) | **Patch** /slb/{name}/loglevel/ | Update loglevel by ID
[**UpdateSlbServerIdByID**](SlbApi.md#UpdateSlbServerIdByID) | **Patch** /slb/{name}/server-id/ | Update server-id by ID


# **CreateSlbByID**
> CreateSlbByID($name, $slb)

Create slb by ID

Create operation of resource: slb


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **slb** | [**Slb**](Slb.md)| slbbody object | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteSlbByID**
> DeleteSlbByID($name)

Delete slb by ID

Delete operation of resource: slb


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbByID**
> Slb ReadSlbByID($name)

Read slb by ID

Read operation of resource: slb


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

[**Slb**](Slb.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbChannelLenByID**
> int32 ReadSlbChannelLenByID($name)

Read channel-len by ID

Read operation of resource: channel-len


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**int32**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbChannelLocByID**
> string ReadSlbChannelLocByID($name)

Read channel-loc by ID

Read operation of resource: channel-loc


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbEgressActionByID**
> string ReadSlbEgressActionByID($name)

Read egress-action by ID

Read operation of resource: egress-action


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbIngressActionByID**
> string ReadSlbIngressActionByID($name)

Read ingress-action by ID

Read operation of resource: ingress-action


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbListByID**
> []Slb ReadSlbListByID()

Read slb by ID

Read operation of resource: slb


### Parameters
This endpoint does not need any parameter.

### Return type

[**[]Slb**](Slb.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbLoglevelByID**
> string ReadSlbLoglevelByID($name)

Read loglevel by ID

Read operation of resource: loglevel


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbServerIdByID**
> int32 ReadSlbServerIdByID($name)

Read server-id by ID

Read operation of resource: server-id


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**int32**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbServiceNameByID**
> string ReadSlbServiceNameByID($name)

Read service-name by ID

Read operation of resource: service-name


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbTypeByID**
> string ReadSlbTypeByID($name)

Read type by ID

Read operation of resource: type


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReadSlbUuidByID**
> string ReadSlbUuidByID($name)

Read uuid by ID

Read operation of resource: uuid


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReplaceSlbByID**
> ReplaceSlbByID($name, $slb)

Replace slb by ID

Replace operation of resource: slb


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **slb** | [**Slb**](Slb.md)| slbbody object | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbByID**
> UpdateSlbByID($name, $slb)

Update slb by ID

Update operation of resource: slb


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **slb** | [**Slb**](Slb.md)| slbbody object | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbChannelLenByID**
> UpdateSlbChannelLenByID($name, $channelLen)

Update channel-len by ID

Update operation of resource: channel-len


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **channelLen** | **int32**| number of bits used for channel | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbChannelLocByID**
> UpdateSlbChannelLocByID($name, $channelLoc)

Update channel-loc by ID

Update operation of resource: channel-loc


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **channelLoc** | **string**| where the channel info located? Default is LSB. | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbEgressActionByID**
> UpdateSlbEgressActionByID($name, $egressAction)

Update egress-action by ID

Update operation of resource: egress-action


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **egressAction** | **string**| Action performed on egress packets | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbIngressActionByID**
> UpdateSlbIngressActionByID($name, $ingressAction)

Update ingress-action by ID

Update operation of resource: ingress-action


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **ingressAction** | **string**| Action performed on ingress packets | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbListByID**
> UpdateSlbListByID($slb)

Update slb by ID

Update operation of resource: slb


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slb** | [**[]Slb**](Slb.md)| slbbody object | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbLoglevelByID**
> UpdateSlbLoglevelByID($name, $loglevel)

Update loglevel by ID

Update operation of resource: loglevel


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **loglevel** | **string**| Defines the logging level of a service instance, from none (OFF) to the most verbose (TRACE) | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateSlbServerIdByID**
> UpdateSlbServerIdByID($name, $serverId)

Update server-id by ID

Update operation of resource: server-id


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| ID of name | 
 **serverId** | **int32**| server id | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

