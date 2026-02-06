// Copyright 2021 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testdata

var (
	// RuleHitReport contains a default and valid rule hit report
	RuleHitReport = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": {
			"reports": [
				{
					"rule_id": "pods_check_containers|POD_CONTAINER_ISSUE",
        			"component": "ccx_rules_ocp.internal.rules.pods_check_containers.report",
        			"type": "rule",
        			"key": "POD_CONTAINER_ISSUE",
        			"details": {
          				"containers": [
            				{
              					"pod": "openshift-apiserver-operator-67667f5cd8-4ctwp",
              					"name": "openshift-apiserver-operator",
              					"ready": false,
              					"restarts": 0
            				},
            				{
              					"pod": "downloads-598f9dd7f9-wwb9g",
              					"name": "download-server",
              					"ready": false,
              					"restarts": 0
            				}
          				],
          				"type": "rule",
          				"error_key": "POD_CONTAINER_ISSUE"
        			},
        			"tags": [],
        			"links": {}
      			}
			]
		}
	}`)

	// FeatureReport contains a default and valid feature report
	FeatureReport = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "cluster_info",
					"component": "fe.features.cluster_info.feature"
				},
				"data": [
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
          				"platform": "AWS",
          				"current_version": "4.6.16",
          				"desired_version": "4.6.16"
					}
				]
			},
			{
				"metadata": {
					"feature_id": "foc",
					"component": "fe.features.foc.feature"
				},
				"data": [
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"operator": "authentication",
						"type": "Degraded",
						"status": true,
						"reason": "RouteStatus_FailedCreate",
						"message": "RouteStatusDegraded: the server is currently unable to handle the request (get routes.route.openshift.io oauth-openshift)",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					},
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"operator": "console",
						"type": "Degraded",
						"status": true,
						"reason": "CustomRouteSync_FailedDeleteCustomRoutes::DefaultRouteSync_FailedDefaultRouteApply::SyncLoopRefresh_InProgress::TrustedCASync_FailedGet",
						"message": "SyncLoopRefreshDegraded: the server is currently unable to handle the request (get routes.route.openshift.io console)\nCustomRouteSyncDegraded: the server is currently unable to handle the request (delete routes.route.openshift.io console-custom)\nDefaultRouteSyncDegraded: the server is currently unable to handle the request (get routes.route.openshift.io console)\nTrustedCASyncDegraded: etcdserver: leader changed",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					}
				]
			},
			{
				"metadata": {
					"feature_id": "alerts",
					"component": "fe.features.alerts.feature"
				},
				"data": [
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "AlertmanagerReceiversNotConfigured",
						"severity": "warning",
						"state": "firing",
						"last_transition_time": "2021-05-05T01:02:36",
						"analysis_time": "2021-05-06T11:39:59",
						"labels": {
							"instance": "",
							"prometheus": "openshift-monitoring/k8s",
							"prometheus_replica": "prometheus-k8s-0"
						},
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202105/05/014446.tar.gz",
						"value": 1
					},
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "Watchdog",
						"severity": "none",
						"state": "firing",
						"last_transition_time": "2021-05-05T01:03:00",
						"analysis_time": "2021-05-06T11:39:59",
						"labels": {
							"instance": "10.20.30.40:24231",
							"prometheus": "openshift-monitoring/k8s",
							"prometheus_replica": "prometheus-k8s-0",
							"type": "elasticsearch"
						},
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202105/05/014446.tar.gz",
						"value": 1
					},
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "Watchdog",
						"severity": "none",
						"state": "firing",
						"last_transition_time": "2021-05-05T01:03:00",
						"analysis_time": "2021-05-06T11:39:59",
						"labels": {
							"cluster_id": "c-morning-dawn-1944",
							"prometheus": "openshift-monitoring/k8s",
							"prometheus_replica": "prometheus-k8s-0",
							"tenant_id": "t-silent-sky-2123"
						},
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202105/05/014446.tar.gz",
						"value": 1
					}
				]
			},
			{
				"metadata": {
					"feature_id": "gatherers",
					"component": "fe.features.gatherer_info.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "name",
							"type": "String"
						},
						{
							"name": "duration_in_ms",
							"type": "Integer"
						},
						{
							"name": "analysis_time",
							"type": "DateTime"
						},
						{
							"name": "path",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "clusterconfig.GatherSAPVsystemIptablesLogs",
						"duration_in_ms": 31,
						"analysis_time": "2021-07-29T15:00:44",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					  },
					  {
						"cluster_id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "clusterconfig.GatherHostSubnet",
						"duration_in_ms": 47,
						"analysis_time": "2021-07-29T15:00:44",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					  }
				]
			},
			{
				"metadata": {
					"feature_id": "available_updates",
					"component": "fe.features.available_updates.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "cluster_version",
							"type": "String"
						},
						{
							"name": "release",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4"
					},
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4"
					}
				]
            },
			{
				"metadata": {
					"feature_id": "conditional_update_conditions",
					"component": "fe.features.conditional_update_conditions.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "cluster_version",
							"type": "String"
						},
						{
							"name": "release",
							"type": "String"
						},
						{
							"name": "recommended",
							"type": "String"
						},
						{
							"name": "reason",
							"type": "String"
						},
						{
							"name": "message",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"recommended": "True",
						"reason": "AsExpected",
						"message": "The update is recommended, because none of the conditional update risks apply to this cluster."
					},
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"recommended": "Unknown",
						"reason": "AsExpected",
						"message": "The update is recommended, because none of the conditional update risks apply to this cluster."
					}
				]
            },
			{
				"metadata": {
					"feature_id": "conditional_update_risks",
					"component": "fe.features.conditional_update_risks.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "cluster_version",
							"type": "String"
						},
						{
							"name": "release",
							"type": "String"
						},
						{
							"name": "risk",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"risk": "AlibabaStorageDriverDemo"
					},
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"risk": "AlibabaStorageDriverDemo"
					}
				]
            }			
		]
	}
	`)

	FeatureReportClusterInfo = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [{
			"metadata": {
				"feature_id": "cluster_info",
				"component": "fe.features.cluster_info.feature"
			},
			"data": [{
				"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"platform": "AWS",
				"current_version": "4.6.16",
				"desired_version": "4.6.16"
			}]
		}]
	}
	`)

	FeatureReportClusterInfoTwoElementsInData = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [{
			"metadata": {
				"feature_id": "cluster_info",
				"component": "fe.features.cluster_info.feature"
			},
			"data": [{
				"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"platform": "AWS",
				"current_version": "4.6.16",
				"desired_version": "4.6.16"
			},
			{
				"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"platform": "AWS",
				"current_version": "4.6.16",
				"desired_version": "4.6.16"
			}]
		}]
	}
	`)

	FeatureReportClusterInfoInvalid = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [{
			"metadata": {
				"feature_id": "cluster_info",
				"component": "fe.features.cluster_info.feature"
			},
			"data": [{
				"networks": "not a struct"
			}]
		}]
	}
	`)

	FeatureReportFOC = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "foc",
					"component": "fe.features.foc.feature"
				},
				"data": [
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"operator": "authentication",
						"type": "Degraded",
						"status": true,
						"reason": "RouteStatus_FailedCreate",
						"message": "RouteStatusDegraded: the server is currently unable to handle the request (get routes.route.openshift.io oauth-openshift)",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					}
				]
			}]
	}
	`)

	FeatureReportFOCInvalid = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "foc",
					"component": "fe.features.foc.feature"
				},
				"data": [
					{
						"status": "not a bool"
					}
				]
			}]
	}
	`)

	FeatureReportAlerts = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [{
			"metadata": {
				"feature_id": "alerts",
				"component": "fe.features.alerts.feature"
			},
			"data": [{
					"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
					"name": "AlertmanagerReceiversNotConfigured",
					"severity": "warning",
					"state": "firing",
					"last_transition_time": "2021-05-05T01:02:36",
					"analysis_time": "2021-05-06T11:39:59",
					"labels": {
						"instance": "",
						"prometheus": "openshift-monitoring/k8s",
						"prometheus_replica": "prometheus-k8s-0"
					},
					"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202105/05/014446.tar.gz",
					"value": 1
				}
			]
		}]
	}
	`)

	FeatureReportAlertsInvalid = []byte(`
	{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [{
			"metadata": {
				"feature_id": "alerts",
				"component": "fe.features.alerts.feature"
			},
			"data": [
				{
					"labels": "not a map"
				}
			]
		}]
	}
	`)

	FeatureReportGatherers = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "gatherers",
					"component": "fe.features.gatherer_info.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "name",
							"type": "String"
						},
						{
							"name": "duration_in_ms",
							"type": "Integer"
						},
						{
							"name": "analysis_time",
							"type": "DateTime"
						},
						{
							"name": "path",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "clusterconfig.GatherSAPVsystemIptablesLogs",
						"duration_in_ms": 31,
						"analysis_time": "2021-07-29T15:00:44",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					  },
					  {
						"cluster_id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"name": "clusterconfig.GatherHostSubnet",
						"duration_in_ms": 47,
						"analysis_time": "2021-07-29T15:00:44",
						"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz"
					  }
				]
			}
		]
	}`)

	FeatureReportGatherersInvalid = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "gatherers",
					"component": "fe.features.gatherer_info.feature"
				},
				"data": [
					{
						"duration_in_ms": "not an int"
					}
				]
			}
		]
	}`)

	FeatureReportWorkloadInfo = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "workload_info_namespaces",
					"component": "fe.features.workload_info_namespaces.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
					  {
						"name": "cluster_id",
						"type": "String"
					  },
					  {
						"name": "namespace",
						"type": "String"
					  },
					  {
						"name": "shapes",
						"type": "List"
					  },
					  {
						"name": "path",
						"type": "String"
					  }
					]
				},
				"data": [
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"namespace": "0LiT6ZNtbpYL",
						"shapes": [
							{
								"restartAlways": true,
								"containers": [
									{
										"imageID": "sha256:53cca031acdd20d285f1a80b5de34df664eee6e36037b9d3521a4df14ed50fd7",
                  						"firstCommand": "N9KxLV2avCo2",
                  						"firstArg": "BuLIUMMJnyP_"
                					}
								],
								"shapeInstances": 6
							},
							{
								"restartAlways": true,
								"initContainers": [
									{
										"imageID": "sha256:9f3c6aa0fa9ecd51f4c110640489f6a074f0e4fd2d950db18249118913f36fd4",
										"firstCommand": "Cl6kTzfbYztA"
								  	},
									{
										"imageID": "sha256:f110a8b30709d390bd57782c270257d9b08fd78b6b4b3a3ad1393e9b018dc06b",
										"firstCommand": "Cl6kTzfbYztA"
								  	},
								  	{
										"imageID": "sha256:0535543d6f48ca5f959fb443b795f5f76c9b019faac3c56e289b71db8f5969eb",
										"firstCommand": "Cl6kTzfbYztA"
								  	},
								  	{
										"imageID": "sha256:392a473b722ff81b443fa9cf290e8f430fb89854e5a0dba559d61adac1a257d1",
										"firstCommand": "Cl6kTzfbYztA"
								  	},
								  	{
										"imageID": "sha256:86a337fcfb2f0538512960ea68f2429336255d3c5f85861f6ec6c279b431f753",
										"firstCommand": "Cl6kTzfbYztA"
								  	},
								  	{
										"imageID": "sha256:86a337fcfb2f0538512960ea68f2429336255d3c5f85861f6ec6c279b431f753",
										"firstCommand": "icTsn2s_EIax"
								  	}
								],
								"containers": [
									{
										"imageID": "sha256:53cca031acdd20d285f1a80b5de34df664eee6e36037b9d3521a4df14ed50fd7",
										"firstCommand": "N9KxLV2avCo2",
										"firstArg": "EbplhSJxzSTF"
								  	}
								],
								"shapeInstances": 6
							}
						]
					}
				]
            }

		]
	}`)

	FeatureReportWorkloadInfoInvalid = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "workload_info_namespaces",
					"component": "fe.features.workload_info_namespaces.feature"
				},
				"data": [
					{
						"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"namespace": "0LiT6ZNtbpYL",
						"shapes": [
							{
								"shapeInstances": "not an int"
							}
						]
					}
				]
            }

		]
	}`)

	ImageLayersReport = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "workload_info_images",
					"component": "fe.features.workload_info_images.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "image_id",
							"type": "String"
						},
						{
							"name": "layers",
							"type": "list"
						},
						{
							"name": "first_command",
							"type": "String"
						},
						{
							"name": "first_arg",
							"type": "String"
						},
						{
							"name": "path",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "ae9900a9-7eb9-4a50-82bd-39ab93177049",
						"image_id": "sha256:05dd0faaf7cf11dedd4efe7eb8574f848094b1604f28a2c3f3f6fbb7b6b2dac9",
						"layers": [
							"sha256:eac1b95df832dc9f172fd1f07e7cb50c1929b118a4249ddd02c6318a677b506a",
							"sha256:47aa3ed2034c4f27622b989b26c06087de17067268a19a1b3642a7e2686cd1a3",
							"sha256:063673bc05f8848a428f045774f9836e1cd29b373f22c1c4c98d94b8d71e77d8",
							"sha256:da314f1f91e1309b4d16ec3e58ac7649574fad70de6a9eef6ecc00ab5261d331",
							"sha256:32658f53f4ca412e4d6dc90f3b4a61ff64a7afb947ddcc8e1f21b8e420a8a778"
						],
						"first_command": "icTsn2s_EIax",
						"first_arg": "2v1NneeWoS_9",
						"analysis_time": "2022-01-25T09:16:39",
						"path": "archives/compressed/ae/ae9900a9-7eb9-4a50-82bd-39ab93177049/202201/25/091559.tar.gz"
					},
					{
						"cluster_id": "ae9900a9-7eb9-4a50-82bd-39ab93177049",
						"image_id": "sha256:08a7fc2587f2e2c2ced90f16cef34582d16aca9c02281350aa35c98d2bc03479",
						"layers": [
							"sha256:eac1b95df832dc9f172fd1f07e7cb50c1929b118a4249ddd02c6318a677b506a",
							"sha256:47aa3ed2034c4f27622b989b26c06087de17067268a19a1b3642a7e2686cd1a3",
							"sha256:063673bc05f8848a428f045774f9836e1cd29b373f22c1c4c98d94b8d71e77d8",
							"sha256:3f5ae93884de9bace5d8bb680e93469ad1cc59984a407be180559fc2482b0192"
						],
						"first_command": "icTsn2s_EIax",
						"first_arg": "2v1NneeWoS_9",
						"analysis_time": "2022-01-25T09:16:39",
						"path": "archives/compressed/ae/ae9900a9-7eb9-4a50-82bd-39ab93177049/202201/25/091559.tar.gz"
					},
					{
						"cluster_id": "ae9900a9-7eb9-4a50-82bd-39ab93177049",
						"image_id": "sha256:3ce8b89d2973dbf10fe15a0843dd6f02d5ff26d2a27d07d2bbb7c23a806a5eca",
						"layers": [
							"sha256:eac1b95df832dc9f172fd1f07e7cb50c1929b118a4249ddd02c6318a677b506a",
							"sha256:47aa3ed2034c4f27622b989b26c06087de17067268a19a1b3642a7e2686cd1a3",
							"sha256:063673bc05f8848a428f045774f9836e1cd29b373f22c1c4c98d94b8d71e77d8",
							"sha256:da314f1f91e1309b4d16ec3e58ac7649574fad70de6a9eef6ecc00ab5261d331",
							"sha256:32efe9ecafb8941f7cd6c44c3b1f1b47eecee53b6ff52aefdd8d73ee9239aef9",
							"sha256:abdd820874221ac4cabaa9c24f905323f42e4deb280cea8f6ab104c0a9732752"
						],
						"first_command": "icTsn2s_EIax",
						"first_arg": "2v1NneeWoS_9",
						"analysis_time": "2022-01-25T09:16:39",
						"path": "archives/compressed/ae/ae9900a9-7eb9-4a50-82bd-39ab93177049/202201/25/091559.tar.gz"
					}
				]
            }
		]
	}`)

	ImageLayersReportInvalid = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "workload_info_images",
					"component": "fe.features.workload_info_images.feature"
				},
				"data": [
					{
						"layers": "not a list"
					}				
				]
            }
		]
	}`)

	AvailableUpdatesReport = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "available_updates",
					"component": "fe.features.available_updates.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "cluster_version",
							"type": "String"
						},
						{
							"name": "release",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4"
					},
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4"
					}
				]
            }
		]
	}`)

	ConditionalUpdateConditionsReport = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "conditional_update_conditions",
					"component": "fe.features.conditional_update_conditions.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "cluster_version",
							"type": "String"
						},
						{
							"name": "release",
							"type": "String"
						},
						{
							"name": "recommended",
							"type": "String"
						},
						{
							"name": "reason",
							"type": "String"
						},
						{
							"name": "message",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"recommended": "True",
						"reason": "AsExpected",
						"message": "The update is recommended, because none of the conditional update risks apply to this cluster."
					},
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"recommended": "Unknown",
						"reason": "AsExpected",
						"message": "The update is recommended, because none of the conditional update risks apply to this cluster."
					}
				]
            }
		]
	}`)

	ConditionalUpdateRisksReport = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "conditional_update_risks",
					"component": "fe.features.conditional_update_risks.feature"
				},
				"schema": {
					"version": "1.0",
					"fields": [
						{
							"name": "cluster_id",
							"type": "String"
						},
						{
							"name": "cluster_version",
							"type": "String"
						},
						{
							"name": "release",
							"type": "String"
						},
						{
							"name": "risk",
							"type": "String"
						}
					]
				},
				"data": [
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"risk": "AlibabaStorageDriverDemo"
					},
					{
						"cluster_id": "2f22800c-52fb-459e-9ba8-2bfd0602ee96",
						"cluster_version": "4.8.12",
						"release": "4.10.0-fc.4",
						"risk": "AlibabaStorageDriverDemo"
					}
				]
            }
		]
	}`)

	AvailableUpdatesReportInvalid = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "available_updates",
					"component": "fe.features.available_updates.feature"
				},
				"data": [
					{
						"cluster_id": true
					}				
				]
            }
		]
	}`)

	ConditionalUpdateConditionsReportInvalid = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "conditional_update_conditions",
					"component": "fe.features.conditional_update_conditions.feature"
				},
				"data": [
					{
						"cluster_id": true
					}				
				]
            }
		]
	}`)

	ConditionalUpdateRisksReportInvalid = []byte(`{
		"path": "archives/compressed/aa/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/202101/20/031044.tar.gz",
		"metadata": {
			"cluster_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"external_organization": "1234567"
		},
		"report": [
			{
				"metadata": {
					"feature_id": "conditional_update_risks",
					"component": "fe.features.conditional_update_risks.feature"
				},
				"data": [
					{
						"cluster_id": true
					}				
				]
            }
		]
	}`)
)
