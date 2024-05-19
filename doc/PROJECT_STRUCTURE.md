.
├── cmd
│   └── main.go
├── config
│   ├── app.dev.yaml
│   ├── app.prod.yaml
│   ├── log.dev.yaml
│   └── log.prod.yaml
├── deployment
│   ├── migration
│   │   ├── 01.profile_down.sql
│   │   ├── 01.profile_up.sql
│   │   └── migrate.sh
│   ├── Dockerfile
│   └── ci_cd.yaml
├── dev_memorianexus
│   ├── dial
│   │   ├── auth.http
│   │   ├── http-client.env.json
│   │   └── profile.http
│   ├── logs
│   │   └── memorianexus.log
│   ├── mysql-data
│   │   ├── #innodb_redo
│   │   │   ├── #ib_redo10_tmp
│   │   │   ├── #ib_redo11_tmp
│   │   │   ├── #ib_redo12_tmp
│   │   │   ├── #ib_redo13_tmp
│   │   │   ├── #ib_redo14_tmp
│   │   │   ├── #ib_redo15_tmp
│   │   │   ├── #ib_redo16_tmp
│   │   │   ├── #ib_redo17_tmp
│   │   │   ├── #ib_redo18_tmp
│   │   │   ├── #ib_redo19_tmp
│   │   │   ├── #ib_redo20_tmp
│   │   │   ├── #ib_redo21_tmp
│   │   │   ├── #ib_redo22_tmp
│   │   │   ├── #ib_redo23_tmp
│   │   │   ├── #ib_redo24_tmp
│   │   │   ├── #ib_redo25_tmp
│   │   │   ├── #ib_redo26_tmp
│   │   │   ├── #ib_redo27_tmp
│   │   │   ├── #ib_redo28_tmp
│   │   │   ├── #ib_redo29_tmp
│   │   │   ├── #ib_redo30_tmp
│   │   │   ├── #ib_redo31_tmp
│   │   │   ├── #ib_redo32_tmp
│   │   │   ├── #ib_redo33_tmp
│   │   │   ├── #ib_redo34_tmp
│   │   │   ├── #ib_redo35_tmp
│   │   │   ├── #ib_redo36_tmp
│   │   │   ├── #ib_redo37_tmp
│   │   │   ├── #ib_redo38_tmp
│   │   │   ├── #ib_redo39_tmp
│   │   │   ├── #ib_redo40_tmp
│   │   │   └── #ib_redo9
│   │   ├── #innodb_temp
│   │   │   ├── temp_1.ibt
│   │   │   ├── temp_10.ibt
│   │   │   ├── temp_2.ibt
│   │   │   ├── temp_3.ibt
│   │   │   ├── temp_4.ibt
│   │   │   ├── temp_5.ibt
│   │   │   ├── temp_6.ibt
│   │   │   ├── temp_7.ibt
│   │   │   ├── temp_8.ibt
│   │   │   └── temp_9.ibt
│   │   ├── memorianexus
│   │   │   ├── profile.ibd
│   │   │   ├── profile_points.ibd
│   │   │   ├── profile_settings_advance.ibd
│   │   │   └── profile_settings_memorization.ibd
│   │   ├── mysql
│   │   │   ├── general_log.CSM
│   │   │   ├── general_log.CSV
│   │   │   ├── general_log_214.sdi
│   │   │   ├── slow_log.CSM
│   │   │   ├── slow_log.CSV
│   │   │   └── slow_log_215.sdi
│   │   ├── performance_schema
│   │   │   ├── accounts_146.sdi
│   │   │   ├── binary_log_trans_190.sdi
│   │   │   ├── cond_instances_81.sdi
│   │   │   ├── data_lock_waits_162.sdi
│   │   │   ├── data_locks_161.sdi
│   │   │   ├── error_log_82.sdi
│   │   │   ├── events_errors_su_140.sdi
│   │   │   ├── events_errors_su_141.sdi
│   │   │   ├── events_errors_su_142.sdi
│   │   │   ├── events_errors_su_143.sdi
│   │   │   ├── events_errors_su_144.sdi
│   │   │   ├── events_stages_cu_112.sdi
│   │   │   ├── events_stages_hi_113.sdi
│   │   │   ├── events_stages_hi_114.sdi
│   │   │   ├── events_stages_su_115.sdi
│   │   │   ├── events_stages_su_116.sdi
│   │   │   ├── events_stages_su_117.sdi
│   │   │   ├── events_stages_su_118.sdi
│   │   │   ├── events_stages_su_119.sdi
│   │   │   ├── events_statement_120.sdi
│   │   │   ├── events_statement_121.sdi
│   │   │   ├── events_statement_122.sdi
│   │   │   ├── events_statement_123.sdi
│   │   │   ├── events_statement_124.sdi
│   │   │   ├── events_statement_125.sdi
│   │   │   ├── events_statement_126.sdi
│   │   │   ├── events_statement_127.sdi
│   │   │   ├── events_statement_128.sdi
│   │   │   ├── events_statement_129.sdi
│   │   │   ├── events_statement_130.sdi
│   │   │   ├── events_statement_131.sdi
│   │   │   ├── events_transacti_132.sdi
│   │   │   ├── events_transacti_133.sdi
│   │   │   ├── events_transacti_134.sdi
│   │   │   ├── events_transacti_135.sdi
│   │   │   ├── events_transacti_136.sdi
│   │   │   ├── events_transacti_137.sdi
│   │   │   ├── events_transacti_138.sdi
│   │   │   ├── events_transacti_139.sdi
│   │   │   ├── events_waits_cur_83.sdi
│   │   │   ├── events_waits_his_84.sdi
│   │   │   ├── events_waits_his_85.sdi
│   │   │   ├── events_waits_sum_86.sdi
│   │   │   ├── events_waits_sum_87.sdi
│   │   │   ├── events_waits_sum_88.sdi
│   │   │   ├── events_waits_sum_89.sdi
│   │   │   ├── events_waits_sum_90.sdi
│   │   │   ├── events_waits_sum_91.sdi
│   │   │   ├── file_instances_92.sdi
│   │   │   ├── file_summary_by__93.sdi
│   │   │   ├── file_summary_by__94.sdi
│   │   │   ├── global_status_182.sdi
│   │   │   ├── global_variables_185.sdi
│   │   │   ├── host_cache_95.sdi
│   │   │   ├── hosts_147.sdi
│   │   │   ├── keyring_componen_192.sdi
│   │   │   ├── keyring_keys_153.sdi
│   │   │   ├── log_status_175.sdi
│   │   │   ├── memory_summary_b_155.sdi
│   │   │   ├── memory_summary_b_156.sdi
│   │   │   ├── memory_summary_b_157.sdi
│   │   │   ├── memory_summary_b_158.sdi
│   │   │   ├── memory_summary_g_154.sdi
│   │   │   ├── metadata_locks_160.sdi
│   │   │   ├── mutex_instances_96.sdi
│   │   │   ├── objects_summary__97.sdi
│   │   │   ├── performance_time_98.sdi
│   │   │   ├── persisted_variab_188.sdi
│   │   │   ├── prepared_stateme_176.sdi
│   │   │   ├── processlist_99.sdi
│   │   │   ├── replication_appl_166.sdi
│   │   │   ├── replication_appl_167.sdi
│   │   │   ├── replication_appl_168.sdi
│   │   │   ├── replication_appl_169.sdi
│   │   │   ├── replication_appl_171.sdi
│   │   │   ├── replication_appl_172.sdi
│   │   │   ├── replication_asyn_173.sdi
│   │   │   ├── replication_asyn_174.sdi
│   │   │   ├── replication_conn_163.sdi
│   │   │   ├── replication_conn_165.sdi
│   │   │   ├── replication_grou_164.sdi
│   │   │   ├── replication_grou_170.sdi
│   │   │   ├── rwlock_instances_100.sdi
│   │   │   ├── session_account__152.sdi
│   │   │   ├── session_connect__151.sdi
│   │   │   ├── session_status_183.sdi
│   │   │   ├── session_variable_186.sdi
│   │   │   ├── setup_actors_101.sdi
│   │   │   ├── setup_consumers_102.sdi
│   │   │   ├── setup_instrument_103.sdi
│   │   │   ├── setup_meters_104.sdi
│   │   │   ├── setup_metrics_105.sdi
│   │   │   ├── setup_objects_106.sdi
│   │   │   ├── setup_threads_107.sdi
│   │   │   ├── socket_instances_148.sdi
│   │   │   ├── socket_summary_b_149.sdi
│   │   │   ├── socket_summary_b_150.sdi
│   │   │   ├── status_by_accoun_178.sdi
│   │   │   ├── status_by_host_179.sdi
│   │   │   ├── status_by_thread_180.sdi
│   │   │   ├── status_by_user_181.sdi
│   │   │   ├── table_handles_159.sdi
│   │   │   ├── table_io_waits_s_108.sdi
│   │   │   ├── table_io_waits_s_109.sdi
│   │   │   ├── table_lock_waits_110.sdi
│   │   │   ├── threads_111.sdi
│   │   │   ├── tls_channel_stat_191.sdi
│   │   │   ├── user_defined_fun_189.sdi
│   │   │   ├── user_variables_b_177.sdi
│   │   │   ├── users_145.sdi
│   │   │   ├── variables_by_thr_184.sdi
│   │   │   └── variables_info_187.sdi
│   │   ├── sys
│   │   │   └── sys_config.ibd
│   │   ├── #ib_16384_0.dblwr
│   │   ├── #ib_16384_1.dblwr
│   │   ├── auto.cnf
│   │   ├── binlog.000001
│   │   ├── binlog.000002
│   │   ├── binlog.000003
│   │   ├── binlog.000004
│   │   ├── binlog.index
│   │   ├── ca-key.pem
│   │   ├── ca.pem
│   │   ├── client-cert.pem
│   │   ├── client-key.pem
│   │   ├── ib_buffer_pool
│   │   ├── ibdata1
│   │   ├── ibtmp1
│   │   ├── mysql.ibd
│   │   ├── mysql.sock -> /var/run/mysqld/mysqld.sock
│   │   ├── mysql_upgrade_history
│   │   ├── private_key.pem
│   │   ├── public_key.pem
│   │   ├── server-cert.pem
│   │   ├── server-key.pem
│   │   ├── undo_001
│   │   └── undo_002
│   ├── redis-data
│   │   └── dump.rdb
│   └── docker-compose.yml
├── doc
│   ├── design
│   │   └── CURVE_USAGE_DESIGN.md
│   ├── API_SPEC.md
│   ├── CODE_STRUCTURE.md
│   ├── CONTRIBUTING.md
│   ├── DEPENDENCIES.md
│   ├── PROJECT_STRUCTURE.md
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── repository
│   │   ├── cache.go
│   │   ├── dc.go
│   │   └── rds.go
│   ├── tests
│   │   ├── integration
│   │   │   └── api_test.go
│   │   └── unit
│   │       └── calculator_test.go
│   └── util
│       ├── idgen.go
│       └── util.go
├── pkg
│   └── memcurve
│       ├── calculator.go
│       ├── calculator_test.go
│       └── curvemodel.go
├── script
│   └── setup_project.sh
├── src
│   ├── core
│   │   ├── analytics
│   │   │   ├── reporter.go
│   │   │   └── types.go
│   │   ├── interfaces
│   │   │   └── port.go
│   │   ├── model
│   │   │   ├── repo.go
│   │   │   └── userfactor.go
│   │   ├── reminder
│   │   │   ├── service.go
│   │   │   └── types.go
│   │   ├── review
│   │   │   ├── scheduler.go
│   │   │   └── session.go
│   │   ├── handlers.go
│   │   └── init.go
│   ├── gw
│   │   ├── error_handler.go
│   │   └── routes.go
│   ├── model
│   │   ├── profile.go
│   │   └── profile_points.go
│   └── module
│       ├── achievement
│       │   └── init.go
│       ├── analytic
│       │   └── init.go
│       ├── book
│       │   └── init.go
│       ├── dungeon
│       │   └── init.go
│       ├── item
│       │   └── init.go
│       ├── nft
│       │   └── init.go
│       ├── operation
│       │   └── init.go
│       ├── profile
│       │   ├── api_getuserpoints.go
│       │   ├── api_getuserprofile.go
│       │   ├── api_getusersettings.go
│       │   ├── api_updateuserprofile.go
│       │   ├── api_updateusersettings.go
│       │   └── init.go
│       ├── system
│       │   └── init.go
│       └── trade
│           └── init.go
├── LICENSE
├── MODULES.puml
├── Makefile
├── README.md
├── go.mod
└── go.sum
