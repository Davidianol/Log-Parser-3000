CREATE TABLE IF NOT EXISTS logs (
                                    id            SERIAL PRIMARY KEY,
                                    filename      TEXT        NOT NULL,
                                    status        TEXT        NOT NULL CHECK (status IN ('processing', 'done', 'error')),
                                    node_count    INT         NOT NULL DEFAULT 0,
                                    uploaded_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    error_message TEXT
);

CREATE TABLE IF NOT EXISTS nodes (
                                     id                SERIAL PRIMARY KEY,
                                     log_id            INT  NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
                                     node_desc         TEXT NOT NULL,
                                     num_ports         INT  NOT NULL,
                                     node_type         INT  NOT NULL CHECK (node_type IN (1, 2, 3)),
                                     system_image_guid TEXT NOT NULL,
                                     node_guid         TEXT NOT NULL,
                                     port_guid         TEXT NOT NULL,
                                     UNIQUE (log_id, node_guid)
);

CREATE TABLE IF NOT EXISTS ports (
                                     id       SERIAL PRIMARY KEY,
                                     node_id  INT  NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
                                     log_id   INT  NOT NULL REFERENCES logs(id)  ON DELETE CASCADE,

                                     node_guid TEXT NOT NULL,
                                     port_guid TEXT NOT NULL,
                                     port_num  INT  NOT NULL,

                                     m_key             TEXT,
                                     gid_prfx          TEXT,
                                     msmlid            INT,
                                     lid               INT,
                                     cap_msk           BIGINT,
                                     m_key_lease_period INT,
                                     diag_code         INT,

                                     link_width_actv   INT,
                                     link_width_sup    INT,
                                     link_width_en     INT,
                                     local_port_num    INT,
                                     link_speed_en     INT,
                                     link_speed_actv   INT,
                                     lmc               INT,
                                     m_key_prot_bits   INT,

                                     link_down_def_state INT,
                                     port_phy_state      INT,
                                     port_state          INT,
                                     link_speed_sup      INT,

                                     vl_arb_high_cap INT,
                                     vl_high_limit   INT,
                                     init_type       INT,
                                     vl_cap          INT,
                                     msmsl           INT,
                                     nmtu            INT,

                                     filter_raw_outb INT,
                                     filter_raw_inb  INT,
                                     part_enf_outb   INT,
                                     part_enf_inb    INT,

                                     op_vls       INT,
                                     ho_q_life    INT,
                                     vl_stall_cnt INT,
                                     mtu_cap      INT,

                                     init_type_reply INT,
                                     vl_arb_low_cap  INT,

                                     p_key_violations INT,
                                     m_key_violations INT,
                                     subn_tmo         INT,

                                     multicast_p_key_trap_suppression_enabled INT,
                                     client_reregister                        INT,
                                     guid_cap                                 INT,
                                     q_key_violations                         INT,
                                     max_credit_hint                          INT,
                                     overrun_errs                             INT,
                                     local_phy_error                          INT,
                                     resp_time_value                          INT,

                                     link_round_trip_latency TEXT,
                                     ooo_sl_mask             TEXT,
                                     cap_msk2                TEXT,
                                     fec_actv                TEXT,
                                     retrans_actv            TEXT
);

CREATE TABLE IF NOT EXISTS nodes_info (
                                          id      SERIAL PRIMARY KEY,
                                          node_id INT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
                                          log_id  INT NOT NULL REFERENCES logs(id)  ON DELETE CASCADE,

    -- GeneralInfo
                                          gi_serial_number TEXT,
                                          gi_part_number   TEXT,
                                          gi_revision      TEXT,
                                          gi_product_name  TEXT,

    -- SwitchInfo
                                          sw_linear_fdb_cap          INT,
                                          sw_random_fdb_cap          INT,
                                          sw_mcast_fdb_cap           INT,
                                          sw_linear_fdb_top          INT,
                                          sw_def_port                INT,
                                          sw_def_mcast_pri_port      INT,
                                          sw_def_mcast_not_pri_port  INT,
                                          sw_life_time_value         INT,
                                          sw_port_state_change       INT,
                                          sw_optimized_slvl_mapping  INT,
                                          sw_lids_per_port           INT,
                                          sw_part_enf_cap            INT,
                                          sw_inb_enf_cap             INT,
                                          sw_outb_enf_cap            INT,
                                          sw_filter_raw_inb_cap      INT,
                                          sw_filter_raw_outb_cap     INT,
                                          sw_enp0                    INT,
                                          sw_mcast_fdb_top           INT,

    -- SharpInfo
                                          sharp_endianness               INT,
                                          sharp_enable_endianness_per_job INT,
                                          sharp_reproducibility_disable  INT
);

CREATE TABLE IF NOT EXISTS topology_groups (
                                               id     SERIAL PRIMARY KEY,
                                               log_id INT  NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
                                               key    TEXT NOT NULL,             -- "switch_active", "host_isolated"
                                               UNIQUE (log_id, key)
);

CREATE TABLE IF NOT EXISTS topology_group_nodes (
                                                    id        SERIAL PRIMARY KEY,
                                                    group_id  INT  NOT NULL REFERENCES topology_groups(id) ON DELETE CASCADE,
                                                    node_guid TEXT NOT NULL
);