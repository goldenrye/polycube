/**
* slb API generated from slb.yang
*
* NOTE: This file is auto generated by polycube-codegen
* https://github.com/polycube-network/polycube-codegen
*/


#pragma once


#include "../base/SlbBase.h"



using namespace polycube::service::model;

class Slb : public SlbBase {
 public:
  Slb(const std::string name, const SlbJsonObject &conf);
  virtual ~Slb();

  void packet_in(polycube::service::Sense sense,
                 polycube::service::PacketInMetadata &md,
                 const std::vector<uint8_t> &packet) override;

  /// <summary>
  /// where the channel info located? Default is LSB.
  /// </summary>
  SlbChannelLocEnum getChannelLoc() override;
  void setChannelLoc(const SlbChannelLocEnum &value) override;

  /// <summary>
  /// number of bits used for channel
  /// </summary>
  uint8_t getChannelLen() override;
  void setChannelLen(const uint8_t &value) override;

  /// <summary>
  /// server id
  /// </summary>
  uint16_t getServerId() override;
  void setServerId(const uint16_t &value) override;

  /// <summary>
  /// Action performed on ingress packets
  /// </summary>
  SlbIngressActionEnum getIngressAction() override;
  void setIngressAction(const SlbIngressActionEnum &value) override;

  /// <summary>
  /// Action performed on egress packets
  /// </summary>
  SlbEgressActionEnum getEgressAction() override;
  void setEgressAction(const SlbEgressActionEnum &value) override;

 private:
  SlbChannelLocEnum ch_loc;
  uint8_t ch_len;
  uint16_t serv_id;
  SlbIngressActionEnum i_act;
  SlbEgressActionEnum e_act;
};
