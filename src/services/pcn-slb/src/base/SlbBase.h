/**
* slb API generated from slb.yang
*
* NOTE: This file is auto generated by polycube-codegen
* https://github.com/polycube-network/polycube-codegen
*/


/* Do not edit this file manually */

/*
* SlbBase.h
*
*
*/

#pragma once

#include "../serializer/SlbJsonObject.h"


#include "polycube/services/transparent_cube.h"



#include "polycube/services/utils.h"
#include "polycube/services/fifo_map.hpp"

#include <spdlog/spdlog.h>

using namespace polycube::service::model;


class SlbBase: public virtual polycube::service::TransparentCube {
 public:
  SlbBase(const std::string name);
  
  virtual ~SlbBase();
  virtual void update(const SlbJsonObject &conf);
  virtual SlbJsonObject toJsonObject();

  /// <summary>
  /// where the channel info located? Default is LSB.
  /// </summary>
  virtual SlbChannelLocEnum getChannelLoc() = 0;
  virtual void setChannelLoc(const SlbChannelLocEnum &value) = 0;

  /// <summary>
  /// number of bits used for channel
  /// </summary>
  virtual uint8_t getChannelLen() = 0;
  virtual void setChannelLen(const uint8_t &value) = 0;

  /// <summary>
  /// server id
  /// </summary>
  virtual uint16_t getServerId() = 0;
  virtual void setServerId(const uint16_t &value) = 0;

  /// <summary>
  /// Action performed on ingress packets
  /// </summary>
  virtual SlbIngressActionEnum getIngressAction() = 0;
  virtual void setIngressAction(const SlbIngressActionEnum &value) = 0;

  /// <summary>
  /// Action performed on egress packets
  /// </summary>
  virtual SlbEgressActionEnum getEgressAction() = 0;
  virtual void setEgressAction(const SlbEgressActionEnum &value) = 0;
};
