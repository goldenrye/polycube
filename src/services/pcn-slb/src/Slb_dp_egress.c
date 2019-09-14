
/*
 * Copyright 2019 The Polycube Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <bcc/helpers.h>

enum {
  SLOWPATH_REASON = 1,
};

enum {
  DROP,      // drop packet
  PASS,      // let packet go
  SLOWPATH,  // send packet to slowpath
  SLB,
};

enum {
    MSB,
    LSB,
};

enum {
    ACTION = 0,
    CH_LOC,
    CH_LEN,
    SERV_ID,
    NUM_PARA,
};

struct serv_const {
    uint32_t serv_value;
    uint32_t msb_mask;
    uint32_t lsb_mask;
    uint8_t serv_bits;
};

BPF_ARRAY(para_map, uint16_t, 4);

static inline
void calc_mask(struct serv_const *serv, uint16_t *para[]) {
    serv->serv_bits = *para[CH_LEN];

    if (*para[CH_LOC] == MSB) {
        serv->serv_value = *para[SERV_ID] << (32-serv->serv_bits);
        serv->msb_mask = ~((1<<(32-serv->serv_bits))-1);
        serv->lsb_mask = (1<<(32-serv->serv_bits))-1;
    } else {
        serv->serv_value = *para[SERV_ID];
        serv->msb_mask = ~((1<<serv->serv_bits)-1);
        serv->lsb_mask = (1<<serv->serv_bits)-1;
    }
}

static int handle_rx(struct CTXTYPE *ctx, struct pkt_metadata *md) {
  unsigned int zero = 0, one = 1, two = 2, three = 3;
  uint16_t action, *para[NUM_PARA];
  struct serv_const serv;

  para[ACTION] = para_map.lookup(&zero);
  if (!para[ACTION]) {
      pcn_log(ctx, LOG_ERR, "parameter %d not set", zero);
      return RX_DROP;
  }
  action = *para[ACTION];

  para[CH_LOC] = para_map.lookup(&one);
  if (!para[CH_LOC]) {
      pcn_log(ctx, LOG_ERR, "parameter %d not set", one);
      return RX_DROP;
  }

  para[CH_LEN] = para_map.lookup(&two);
  if (!para[CH_LEN]) {
      pcn_log(ctx, LOG_ERR, "parameter %d not set", two);
      return RX_DROP;
  }

  para[SERV_ID] = para_map.lookup(&three);
  if (!para[SERV_ID]) {
      pcn_log(ctx, LOG_ERR, "parameter %d not set", three);
      return RX_DROP;
  }

  // what action should be performed in the packet?
  switch (action) {
  case DROP:
    pcn_log(ctx, LOG_DEBUG, "Slb egress: dropping packet");
    return RX_DROP;
  case PASS:
    pcn_log(ctx, LOG_DEBUG, "Slb egress: passing packet");
    return RX_OK;
  case SLOWPATH:
    pcn_log(ctx, LOG_DEBUG, "Slb egress: sending packet to slow path");
    return pcn_pkt_controller(ctx, md, SLOWPATH_REASON);
  case SLB:
    pcn_log(ctx, LOG_DEBUG, "Slb egress: slb process");
    calc_mask(&serv, para);
    return RX_OK;
  default:
    pcn_log(ctx, LOG_ERR, "Slb egress: bad action %d", action);
    return RX_DROP;
  }

  return RX_OK;
}
