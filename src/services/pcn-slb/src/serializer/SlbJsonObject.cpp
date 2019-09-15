/**
* slb API generated from slb.yang
*
* NOTE: This file is auto generated by polycube-codegen
* https://github.com/polycube-network/polycube-codegen
*/


/* Do not edit this file manually */



#include "SlbJsonObject.h"
#include <regex>

namespace polycube {
namespace service {
namespace model {

SlbJsonObject::SlbJsonObject() {
  m_nameIsSet = false;
  m_channelLoc = SlbChannelLocEnum::LSB;
  m_channelLocIsSet = true;
  m_channelLenIsSet = false;
  m_serverIdIsSet = false;
  m_ingressAction = SlbIngressActionEnum::PASS;
  m_ingressActionIsSet = true;
  m_egressAction = SlbEgressActionEnum::PASS;
  m_egressActionIsSet = true;
}

SlbJsonObject::SlbJsonObject(const nlohmann::json &val) :
  JsonObjectBase(val) {
  m_nameIsSet = false;
  m_channelLocIsSet = false;
  m_channelLenIsSet = false;
  m_serverIdIsSet = false;
  m_ingressActionIsSet = false;
  m_egressActionIsSet = false;


  if (val.count("name")) {
    setName(val.at("name").get<std::string>());
  }

  if (val.count("channel-loc")) {
    setChannelLoc(string_to_SlbChannelLocEnum(val.at("channel-loc").get<std::string>()));
  }

  if (val.count("channel-len")) {
    setChannelLen(val.at("channel-len").get<uint8_t>());
  }

  if (val.count("server-id")) {
    setServerId(val.at("server-id").get<uint16_t>());
  }

  if (val.count("ingress-action")) {
    setIngressAction(string_to_SlbIngressActionEnum(val.at("ingress-action").get<std::string>()));
  }

  if (val.count("egress-action")) {
    setEgressAction(string_to_SlbEgressActionEnum(val.at("egress-action").get<std::string>()));
  }
}

nlohmann::json SlbJsonObject::toJson() const {
  nlohmann::json val = nlohmann::json::object();
  if (!getBase().is_null()) {
    val.update(getBase());
  }

  if (m_nameIsSet) {
    val["name"] = m_name;
  }

  if (m_channelLocIsSet) {
    val["channel-loc"] = SlbChannelLocEnum_to_string(m_channelLoc);
  }

  if (m_channelLenIsSet) {
    val["channel-len"] = m_channelLen;
  }

  if (m_serverIdIsSet) {
    val["server-id"] = m_serverId;
  }

  if (m_ingressActionIsSet) {
    val["ingress-action"] = SlbIngressActionEnum_to_string(m_ingressAction);
  }

  if (m_egressActionIsSet) {
    val["egress-action"] = SlbEgressActionEnum_to_string(m_egressAction);
  }

  return val;
}

std::string SlbJsonObject::getName() const {
  return m_name;
}

void SlbJsonObject::setName(std::string value) {
  m_name = value;
  m_nameIsSet = true;
}

bool SlbJsonObject::nameIsSet() const {
  return m_nameIsSet;
}



SlbChannelLocEnum SlbJsonObject::getChannelLoc() const {
  return m_channelLoc;
}

void SlbJsonObject::setChannelLoc(SlbChannelLocEnum value) {
  m_channelLoc = value;
  m_channelLocIsSet = true;
}

bool SlbJsonObject::channelLocIsSet() const {
  return m_channelLocIsSet;
}

void SlbJsonObject::unsetChannelLoc() {
  m_channelLocIsSet = false;
}

std::string SlbJsonObject::SlbChannelLocEnum_to_string(const SlbChannelLocEnum &value){
  switch(value) {
    case SlbChannelLocEnum::MSB:
      return std::string("msb");
    case SlbChannelLocEnum::LSB:
      return std::string("lsb");
    default:
      throw std::runtime_error("Bad Slb channelLoc");
  }
}

SlbChannelLocEnum SlbJsonObject::string_to_SlbChannelLocEnum(const std::string &str){
  if (JsonObjectBase::iequals("msb", str))
    return SlbChannelLocEnum::MSB;
  if (JsonObjectBase::iequals("lsb", str))
    return SlbChannelLocEnum::LSB;
  throw std::runtime_error("Slb channelLoc is invalid");
}
uint8_t SlbJsonObject::getChannelLen() const {
  return m_channelLen;
}

void SlbJsonObject::setChannelLen(uint8_t value) {
  m_channelLen = value;
  m_channelLenIsSet = true;
}

bool SlbJsonObject::channelLenIsSet() const {
  return m_channelLenIsSet;
}

void SlbJsonObject::unsetChannelLen() {
  m_channelLenIsSet = false;
}

uint16_t SlbJsonObject::getServerId() const {
  return m_serverId;
}

void SlbJsonObject::setServerId(uint16_t value) {
  m_serverId = value;
  m_serverIdIsSet = true;
}

bool SlbJsonObject::serverIdIsSet() const {
  return m_serverIdIsSet;
}

void SlbJsonObject::unsetServerId() {
  m_serverIdIsSet = false;
}

SlbIngressActionEnum SlbJsonObject::getIngressAction() const {
  return m_ingressAction;
}

void SlbJsonObject::setIngressAction(SlbIngressActionEnum value) {
  m_ingressAction = value;
  m_ingressActionIsSet = true;
}

bool SlbJsonObject::ingressActionIsSet() const {
  return m_ingressActionIsSet;
}

void SlbJsonObject::unsetIngressAction() {
  m_ingressActionIsSet = false;
}

std::string SlbJsonObject::SlbIngressActionEnum_to_string(const SlbIngressActionEnum &value){
  switch(value) {
    case SlbIngressActionEnum::DROP:
      return std::string("drop");
    case SlbIngressActionEnum::PASS:
      return std::string("pass");
    case SlbIngressActionEnum::SLOWPATH:
      return std::string("slowpath");
    default:
      throw std::runtime_error("Bad Slb ingressAction");
  }
}

SlbIngressActionEnum SlbJsonObject::string_to_SlbIngressActionEnum(const std::string &str){
  if (JsonObjectBase::iequals("drop", str))
    return SlbIngressActionEnum::DROP;
  if (JsonObjectBase::iequals("pass", str))
    return SlbIngressActionEnum::PASS;
  if (JsonObjectBase::iequals("slowpath", str))
    return SlbIngressActionEnum::SLOWPATH;
  throw std::runtime_error("Slb ingressAction is invalid");
}
SlbEgressActionEnum SlbJsonObject::getEgressAction() const {
  return m_egressAction;
}

void SlbJsonObject::setEgressAction(SlbEgressActionEnum value) {
  m_egressAction = value;
  m_egressActionIsSet = true;
}

bool SlbJsonObject::egressActionIsSet() const {
  return m_egressActionIsSet;
}

void SlbJsonObject::unsetEgressAction() {
  m_egressActionIsSet = false;
}

std::string SlbJsonObject::SlbEgressActionEnum_to_string(const SlbEgressActionEnum &value){
  switch(value) {
    case SlbEgressActionEnum::DROP:
      return std::string("drop");
    case SlbEgressActionEnum::PASS:
      return std::string("pass");
    case SlbEgressActionEnum::SLOWPATH:
      return std::string("slowpath");
    case SlbEgressActionEnum::SLB:
      return std::string("slb");
    default:
      throw std::runtime_error("Bad Slb egressAction");
  }
}

SlbEgressActionEnum SlbJsonObject::string_to_SlbEgressActionEnum(const std::string &str){
  if (JsonObjectBase::iequals("drop", str))
    return SlbEgressActionEnum::DROP;
  if (JsonObjectBase::iequals("pass", str))
    return SlbEgressActionEnum::PASS;
  if (JsonObjectBase::iequals("slowpath", str))
    return SlbEgressActionEnum::SLOWPATH;
  if (JsonObjectBase::iequals("slb", str))
    return SlbEgressActionEnum::SLB;
  throw std::runtime_error("Slb egressAction is invalid");
}

}
}
}

