/**
 * bridge API generated from bridge.yang
 *
 * NOTE: This file is auto generated by polycube-codegen
 * https://github.com/polycube-network/polycube-codegen
 */

/* Do not edit this file manually */

#include "PortsAccessBase.h"
#include "../Bridge.h"

PortsAccessBase::PortsAccessBase(Ports &parent) : parent_(parent) {}

PortsAccessBase::~PortsAccessBase() {}

void PortsAccessBase::update(const PortsAccessJsonObject &conf) {
  if (conf.vlanidIsSet()) {
    setVlanid(conf.getVlanid());
  }
}

PortsAccessJsonObject PortsAccessBase::toJsonObject() {
  PortsAccessJsonObject conf;

  conf.setVlanid(getVlanid());

  return conf;
}

std::shared_ptr<spdlog::logger> PortsAccessBase::logger() {
  return parent_.logger();
}