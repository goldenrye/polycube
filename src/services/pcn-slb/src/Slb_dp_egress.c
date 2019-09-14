
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
#include <bcc/proto.h>
#include <uapi/linux/in.h>
#include <uapi/linux/ip.h>
#include <uapi/linux/tcp.h>
#include <uapi/linux/pkt_cls.h>

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
    uint8_t  serv_bits;
    uint8_t  ch_loc;
};

struct event_t {
    u32 sip;
    u32 dip;
    u16 sport;
    u16 dport;
    u32 ts_val;
    u32 ts_ecr;
    u32 ts_val_orig;
    u32 ts_xsb;
};

struct ts_info {
    u32 *ts_val_orig;
    u32 *ts_val;
    u32 *ts_ecr;
    u32 *ts_xsb;
};

struct sess_key {
    u32 sip;
    u32 dip;
    u16 sport;
    u16 dport;
};

BPF_ARRAY(para_map, uint16_t, 4);
BPF_TABLE("array", int, u64, counter, 1);
BPF_TABLE("array", int, u64, ts_counter, 1);
BPF_TABLE_SHARED("hash", struct sess_key, u32, sess2ts, 1024);
BPF_PERF_OUTPUT(events);

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
    serv->ch_loc = *para[CH_LOC];
}

static inline
int x_ts_2_sid(void *ptr, u8 tcp_option_len, void *data_end,
                struct sess_key *sess, struct ts_info *ts_info,
                struct serv_const *serv) {
    u8 remain_len = tcp_option_len;
#pragma unroll
    for (u8 i=0; i<5; i++) {
        if (ptr + 10 > data_end || remain_len <= 0)
            break;

        u8 kind = *(u8*)ptr;
        if (kind == 0) {
            ptr++;
            remain_len--;
            break;
        }

        if (kind == 1) {
            ptr++;
            remain_len--;
            continue;
        }

        u8 len = *(u8*)(ptr+1);

        if (kind == 8) {
            u32 xsb;
            *ts_info->ts_val_orig = *ts_info->ts_val = *(u32*)(ptr+2);
            if (serv->ch_loc == MSB)
                xsb = (ntohl(*ts_info->ts_val) & serv->msb_mask) >> (32-serv->serv_bits);
            else
                xsb = ntohl(*ts_info->ts_val) & serv->lsb_mask;

            u32 *val = sess2ts.lookup_or_init(sess, &xsb);
            if (*val != xsb) {
                sess2ts.update(sess, &xsb);
                int zero = 0;
                u64 *value = ts_counter.lookup(&zero);
                if (value) {
                    *value += 1;
                }
            }
            *ts_info->ts_xsb = htonl(xsb);

            // translate ts_val
            if (serv->ch_loc == MSB)
                *ts_info->ts_val = htonl((ntohl(*ts_info->ts_val) & serv->lsb_mask) + serv->serv_value);
            else
                *ts_info->ts_val = htonl((ntohl(*ts_info->ts_val) & serv->msb_mask) + serv->serv_value);

            *(u32*)(ptr+2) = *ts_info->ts_val;
            *ts_info->ts_ecr = *(u32*)(ptr+6);

            break;
        }
        ptr += len;
        remain_len -= len;
    }
    return 0;
}

static inline
int slb_egress_handler(struct CTXTYPE *skb, struct serv_const *serv) {
    int zero = 0;
    u64 *value;
    struct event_t event_data = {};

    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    struct ethhdr *eth = data;
    struct iphdr *iph = data + sizeof(*eth);
    struct tcphdr *tcph;
    void *ptr;

    if (data + sizeof(*eth) + sizeof(*iph) > data_end)
        return RX_OK;

    value = counter.lookup(&zero);
    if (value) {
        *value += 1;
    }

    if (eth->h_proto != htons(ETH_P_IP) || iph->protocol != IPPROTO_TCP) {
        return RX_OK;
    }
    ptr = data + sizeof(*eth) + sizeof(*iph);

    u8 ip_len = iph->ihl<<2;
    u8 extra_len = ip_len - sizeof(*iph);
    ptr += extra_len;

    if (ptr + sizeof(*tcph) > data_end)
        return RX_OK;

    tcph = ptr;
    ptr += sizeof(*tcph);

    u32 ts_val = 0, ts_val_orig = 0;
    u32 ts_ecr = 0;
    u32 ts_xsb = 0;
    if (tcph->doff > 5) {
        u8 tcp_len = tcph->doff<<2;
        u8 tcp_option_len = tcp_len - sizeof(*tcph);

        if (tcp_option_len < 4 || ptr + tcp_option_len > data_end)
            return RX_OK;

        struct ts_info tsi;
        tsi.ts_val = &ts_val;
        tsi.ts_ecr = &ts_ecr;
        tsi.ts_val_orig = &ts_val_orig;
        tsi.ts_xsb = &ts_xsb;

        struct sess_key sess;
        sess.sip = iph->saddr;
        sess.dip = iph->daddr;
        sess.sport = tcph->source;
        sess.dport = tcph->dest;
        x_ts_2_sid(ptr, tcp_option_len, data_end, &sess, &tsi, serv);
    }

    event_data.sip = iph->saddr;
    event_data.dip = iph->daddr;
    event_data.sport = tcph->source;
    event_data.dport = tcph->dest;
    event_data.ts_val = ts_val;
    event_data.ts_ecr = ts_ecr;
    event_data.ts_val_orig = ts_val_orig;
    event_data.ts_xsb = ts_xsb;
    events.perf_submit_skb(skb, skb->len, &event_data, sizeof(event_data));
    return RX_OK;
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
    return slb_egress_handler(ctx, &serv);
  default:
    pcn_log(ctx, LOG_ERR, "Slb egress: bad action %d", action);
    return RX_DROP;
  }

  return RX_OK;
}
