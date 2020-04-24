/**
* dynmon API generated from dynmon.yang
*
* NOTE: This file is auto generated by polycube-codegen
* https://github.com/polycube-network/polycube-codegen
*/


/* Do not edit this file manually */

/*
* DynmonApi.h
*
*/

#pragma once

#define POLYCUBE_SERVICE_NAME "dynmon"


#include "polycube/services/response.h"
#include "polycube/services/shared_lib_elements.h"

#include "DataplaneJsonObject.h"
#include "DataplaneMetricsJsonObject.h"
#include "DataplaneMetricsOpenMetricsMetadataJsonObject.h"
#include "DataplaneMetricsOpenMetricsMetadataLabelsJsonObject.h"
#include "DynmonJsonObject.h"
#include "MetricsJsonObject.h"
#include <vector>


#ifdef __cplusplus
extern "C" {
#endif

Response create_dynmon_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);
Response create_dynmon_dataplane_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);
Response delete_dynmon_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response delete_dynmon_dataplane_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_code_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_list_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_map_name_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_open_metrics_metadata_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_open_metrics_metadata_help_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_open_metrics_metadata_labels_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_open_metrics_metadata_labels_list_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_open_metrics_metadata_labels_value_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_metrics_open_metrics_metadata_type_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_dataplane_name_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_list_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_metrics_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_metrics_list_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_metrics_timestamp_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_metrics_value_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response read_dynmon_open_metrics_by_id_handler(const char *name, const Key *keys, size_t num_keys);
Response replace_dynmon_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);
Response replace_dynmon_dataplane_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);
Response update_dynmon_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);
Response update_dynmon_dataplane_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);
Response update_dynmon_list_by_id_handler(const char *name, const Key *keys, size_t num_keys, const char *value);

Response dynmon_dataplane_metrics_list_by_id_help(const char *name, const Key *keys, size_t num_keys);
Response dynmon_dataplane_metrics_open_metrics_metadata_labels_list_by_id_help(const char *name, const Key *keys, size_t num_keys);
Response dynmon_list_by_id_help(const char *name, const Key *keys, size_t num_keys);
Response dynmon_metrics_list_by_id_help(const char *name, const Key *keys, size_t num_keys);


#ifdef __cplusplus
}
#endif

