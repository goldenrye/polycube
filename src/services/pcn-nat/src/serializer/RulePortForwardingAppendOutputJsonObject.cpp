/**
* nat API
* nat API generated from nat.yang
*
* OpenAPI spec version: 1.0.0
*
* NOTE: This class is auto generated by the swagger code generator program.
* https://github.com/polycube-network/swagger-codegen.git
* branch polycube
*/


/* Do not edit this file manually */



#include "RulePortForwardingAppendOutputJsonObject.h"
#include <regex>

namespace io {
namespace swagger {
namespace server {
namespace model {

RulePortForwardingAppendOutputJsonObject::RulePortForwardingAppendOutputJsonObject() {
  m_idIsSet = false;
}

RulePortForwardingAppendOutputJsonObject::RulePortForwardingAppendOutputJsonObject(const nlohmann::json &val) :
  JsonObjectBase(val) {
  m_idIsSet = false;


  if (val.count("id")) {
    setId(val.at("id").get<uint32_t>());
  }
}

nlohmann::json RulePortForwardingAppendOutputJsonObject::toJson() const {
  nlohmann::json val = nlohmann::json::object();
  if (!getBase().is_null()) {
    val.update(getBase());
  }

  if (m_idIsSet) {
    val["id"] = m_id;
  }

  return val;
}

uint32_t RulePortForwardingAppendOutputJsonObject::getId() const {
  return m_id;
}

void RulePortForwardingAppendOutputJsonObject::setId(uint32_t value) {
  m_id = value;
  m_idIsSet = true;
}

bool RulePortForwardingAppendOutputJsonObject::idIsSet() const {
  return m_idIsSet;
}

void RulePortForwardingAppendOutputJsonObject::unsetId() {
  m_idIsSet = false;
}


}
}
}
}

