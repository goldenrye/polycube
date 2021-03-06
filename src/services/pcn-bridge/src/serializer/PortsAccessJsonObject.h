/**
* bridge API generated from bridge.yang
*
* NOTE: This file is auto generated by polycube-codegen
* https://github.com/polycube-network/polycube-codegen
*/


/* Do not edit this file manually */

/*
* PortsAccessJsonObject.h
*
*
*/

#pragma once


#include "JsonObjectBase.h"


namespace polycube {
namespace service {
namespace model {


/// <summary>
///
/// </summary>
class  PortsAccessJsonObject : public JsonObjectBase {
public:
  PortsAccessJsonObject();
  PortsAccessJsonObject(const nlohmann::json &json);
  ~PortsAccessJsonObject() final = default;
  nlohmann::json toJson() const final;


  /// <summary>
  /// VLAN associated with this interface
  /// </summary>
  uint16_t getVlanid() const;
  void setVlanid(uint16_t value);
  bool vlanidIsSet() const;
  void unsetVlanid();

private:
  uint16_t m_vlanid;
  bool m_vlanidIsSet;
};

}
}
}

