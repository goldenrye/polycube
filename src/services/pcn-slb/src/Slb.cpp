/**
* slb API generated from slb.yang
*
* NOTE: This file is auto generated by polycube-codegen
* https://github.com/polycube-network/polycube-codegen
*/


// TODO: Modify these methods with your own implementation


#include "Slb.h"
#include "Slb_dp_ingress.h"
#include "Slb_dp_egress.h"

using namespace polycube::service;

enum {
    ACTION = 0,
    CH_LOC,
    CH_LEN,
    SERV_ID,
    NUM_PARA,
};

Slb::Slb(const std::string name, const SlbJsonObject &conf)
  : TransparentCube(conf.getBase(), { slb_code_ingress }, { slb_code_egress }),
    SlbBase(name) {
  logger()->info("Creating Slb instance");
  if (conf.channelLocIsSet()) {
    setChannelLoc(conf.getChannelLoc());
  }
  if (conf.channelLenIsSet()) {
    setChannelLen(conf.getChannelLen());
  }
  if (conf.serverIdIsSet()) {
     setServerId(conf.getServerId());
  }
    setIngressAction(conf.getIngressAction());
    setEgressAction(conf.getEgressAction());
}


Slb::~Slb() {
  logger()->info("Destroying Slb instance");
}

void Slb::packet_in(polycube::service::Sense sense,
    polycube::service::PacketInMetadata &md,
    const std::vector<uint8_t> &packet) {
    logger()->debug("Packet received");
}

SlbChannelLocEnum Slb::getChannelLoc() {
    return ch_loc;
}

void Slb::setChannelLoc(const SlbChannelLocEnum &value) {
    ch_loc = value;
    uint8_t loc = static_cast<uint8_t>(value);
    auto t = get_array_table<uint16_t>("para_map", 0, ProgramType::EGRESS);
    t.set(CH_LOC, loc);
}

uint8_t Slb::getChannelLen() {
    return ch_len;
}

void Slb::setChannelLen(const uint8_t &value) {
    ch_len = value;
    uint8_t len = static_cast<uint8_t>(value);
    auto t = get_array_table<uint16_t>("para_map", 0, ProgramType::EGRESS);
    t.set(CH_LEN, len);
}

uint16_t Slb::getServerId() {
    return serv_id;
}

void Slb::setServerId(const uint16_t &value) {
    serv_id = value;
    uint16_t sid = static_cast<uint16_t>(value);
    auto t = get_array_table<uint16_t>("para_map", 0, ProgramType::EGRESS);
    t.set(SERV_ID, sid);
}

SlbIngressActionEnum Slb::getIngressAction() {
    return i_act;
}

void Slb::setIngressAction(const SlbIngressActionEnum &value) {
    i_act = value;
    uint8_t action = static_cast<uint8_t>(value);
    auto t = get_array_table<uint8_t>("action_map", 0, ProgramType::INGRESS);
    t.set(0x0, action);
}

SlbEgressActionEnum Slb::getEgressAction() {
    return e_act;
}

void Slb::setEgressAction(const SlbEgressActionEnum &value) {
    e_act = value;
    uint8_t action = static_cast<uint8_t>(value);
    auto t = get_array_table<uint16_t>("para_map", 0, ProgramType::EGRESS);
    t.set(ACTION, action);
}


