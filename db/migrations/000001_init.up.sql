DROP TABLE IF EXISTS booking_payment CASCADE;
DROP TABLE IF EXISTS lead CASCADE;
DROP TABLE IF EXISTS role CASCADE;
DROP TABLE IF EXISTS tenant_addon_purchase CASCADE;
DROP TABLE IF EXISTS tag_assigned CASCADE;
DROP TABLE IF EXISTS shared_content_log CASCADE;
DROP TABLE IF EXISTS mst_plan_price CASCADE;
DROP TABLE IF EXISTS user_session CASCADE;
DROP TABLE IF EXISTS tenant_setting CASCADE;
DROP TABLE IF EXISTS message_log CASCADE;
DROP TABLE IF EXISTS user_workspace CASCADE;
DROP TABLE IF EXISTS workspace CASCADE;
DROP TABLE IF EXISTS workspace_assignment_rule CASCADE;
DROP TABLE IF EXISTS mst_permission CASCADE;
DROP TABLE IF EXISTS tenant_number_series CASCADE;
DROP TABLE IF EXISTS mst_universal CASCADE;
DROP TABLE IF EXISTS document CASCADE;
DROP TABLE IF EXISTS user_device CASCADE;
DROP TABLE IF EXISTS campaign_import CASCADE;
DROP TABLE IF EXISTS call_log CASCADE;
DROP TABLE IF EXISTS booking_cancellation CASCADE;
DROP TABLE IF EXISTS campaign_member CASCADE;
DROP TABLE IF EXISTS role_permission CASCADE;
DROP TABLE IF EXISTS unit CASCADE;
DROP TABLE IF EXISTS calling_list_record CASCADE;
DROP TABLE IF EXISTS user_tenant CASCADE;
DROP TABLE IF EXISTS booking_parking CASCADE;
DROP TABLE IF EXISTS unit_hold CASCADE;
DROP TABLE IF EXISTS tenant_subscription_payment CASCADE;
DROP TABLE IF EXISTS property_block_stage_progress CASCADE;
DROP TABLE IF EXISTS parking CASCADE;
DROP TABLE IF EXISTS mst_module CASCADE;
DROP TABLE IF EXISTS note CASCADE;
DROP TABLE IF EXISTS quotation CASCADE;
DROP TABLE IF EXISTS tenant CASCADE;
DROP TABLE IF EXISTS mst_status CASCADE;
DROP TABLE IF EXISTS booking_charge CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS activity CASCADE;
DROP TABLE IF EXISTS mst_currency CASCADE;
DROP TABLE IF EXISTS property CASCADE;
DROP TABLE IF EXISTS tag CASCADE;
DROP TABLE IF EXISTS booking CASCADE;
DROP TABLE IF EXISTS calling_list CASCADE;
DROP TABLE IF EXISTS campaign CASCADE;
DROP TABLE IF EXISTS message_template CASCADE;
DROP TABLE IF EXISTS user_auth CASCADE;
DROP TABLE IF EXISTS opportunity CASCADE;
DROP TABLE IF EXISTS property_block CASCADE;
DROP TABLE IF EXISTS booking_storage CASCADE;
DROP TABLE IF EXISTS property_transaction CASCADE;
DROP TABLE IF EXISTS commission CASCADE;
DROP TABLE IF EXISTS tenant_subscription CASCADE;
DROP TABLE IF EXISTS commission_payment CASCADE;
DROP TABLE IF EXISTS tenant_usage_ledger CASCADE;
DROP TABLE IF EXISTS user_otp CASCADE;
DROP TABLE IF EXISTS booking_discount CASCADE;
DROP TABLE IF EXISTS mst_master_type CASCADE;
DROP TABLE IF EXISTS audit_log CASCADE;
DROP TABLE IF EXISTS webhook_log CASCADE;
DROP TABLE IF EXISTS saved_filter CASCADE;
DROP TABLE IF EXISTS contact CASCADE;
DROP TABLE IF EXISTS storage CASCADE;
DROP TABLE IF EXISTS user_login_history CASCADE;
DROP TABLE IF EXISTS quotation_charge CASCADE;
DROP TABLE IF EXISTS media CASCADE;
DROP TABLE IF EXISTS site_visit CASCADE;
DROP TABLE IF EXISTS booking_refund CASCADE;
DROP TABLE IF EXISTS message_queue CASCADE;
DROP TABLE IF EXISTS mst_plan CASCADE;
DROP TABLE IF EXISTS quotation_discount CASCADE;

-- PostgreSQL Database Dump Migration Script
-- Target Engine: PostgreSQL 18 (Optimized for UUIDv7 Application Strategies)

-- --------------------------------------------------------
-- Table structure for table `activity`
-- --------------------------------------------------------
CREATE TABLE activity (
  activity_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  opportunity_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  activity_type_universal_uuid UUID DEFAULT NULL,
  status_uuid UUID NOT NULL,
  priority_universal_uuid UUID DEFAULT NULL,
  activity_datetime TIMESTAMPTZ NOT NULL,
  followup_datetime TIMESTAMPTZ DEFAULT NULL,
  performed_by_user_uuid UUID DEFAULT NULL,
  remarks VARCHAR(2000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (activity_uuid)
);

CREATE INDEX idx_activity_tenant_uuid ON activity (tenant_uuid);
CREATE INDEX idx_activity_workspace_uuid ON activity (workspace_uuid);
CREATE INDEX idx_activity_opportunity_uuid ON activity (opportunity_uuid);
CREATE INDEX idx_activity_followup_datetime ON activity (followup_datetime);
CREATE INDEX idx_activity_performed_by_user_uuid ON activity (performed_by_user_uuid);
CREATE INDEX idx_activity_is_deleted ON activity (is_deleted);
CREATE INDEX idx_activity_module_entity ON activity (module_uuid, entity_uuid);
CREATE INDEX idx_activity_module_entity_created_at ON activity (module_uuid, entity_uuid, created_at);
CREATE INDEX idx_activity_priority_universal_uuid ON activity (priority_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `audit_log`
-- --------------------------------------------------------
CREATE TABLE audit_log (
  audit_log_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  action_name VARCHAR(100) NOT NULL,
  old_value_json JSONB DEFAULT NULL,
  new_value_json JSONB DEFAULT NULL,
  action_by_user_uuid UUID DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  PRIMARY KEY (audit_log_uuid)
);

CREATE INDEX idx_audit_log_entity_uuid ON audit_log (entity_uuid);
CREATE INDEX idx_audit_log_created_at ON audit_log (created_at);
CREATE INDEX idx_audit_log_module_entity ON audit_log (module_uuid, entity_uuid);
CREATE INDEX idx_audit_log_module_entity_created_at ON audit_log (module_uuid, entity_uuid, created_at);

-- --------------------------------------------------------
-- Table structure for table `booking`
-- --------------------------------------------------------
CREATE TABLE booking (
  booking_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  status_uuid UUID NOT NULL,
  opportunity_uuid UUID NOT NULL,
  quotation_uuid UUID DEFAULT NULL,
  contact_uuid UUID NOT NULL,
  unit_uuid UUID NOT NULL,
  payment_plan_universal_uuid UUID DEFAULT NULL,
  booking_number VARCHAR(100) NOT NULL,
  booking_date DATE NOT NULL,
  agreement_value DECIMAL(18,2) DEFAULT NULL,
  received_amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  refunded_amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  balance_amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  booking_amount DECIMAL(18,2) DEFAULT NULL,
  loan_amount DECIMAL(18,2) DEFAULT NULL,
  loan_bank_name VARCHAR(255) DEFAULT NULL,
  co_applicant_details JSONB DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_uuid),
  CONSTRAINT uk_workspace_booking_number UNIQUE (workspace_uuid, booking_number),
  CONSTRAINT uk_booking_unit_uuid UNIQUE (unit_uuid)
);

CREATE INDEX idx_booking_tenant_uuid ON booking (tenant_uuid);
CREATE INDEX idx_booking_workspace_uuid ON booking (workspace_uuid);
CREATE INDEX idx_booking_status_uuid ON booking (status_uuid);
CREATE INDEX idx_booking_opportunity_uuid ON booking (opportunity_uuid);
CREATE INDEX idx_booking_contact_uuid ON booking (contact_uuid);
CREATE INDEX idx_booking_unit_uuid ON booking (unit_uuid);
CREATE INDEX idx_booking_number ON booking (booking_number);
CREATE INDEX idx_booking_date ON booking (booking_date);
CREATE INDEX idx_booking_is_deleted ON booking (is_deleted);
CREATE INDEX idx_booking_quotation_uuid ON booking (quotation_uuid);

-- --------------------------------------------------------
-- Table structure for table `booking_cancellation`
-- --------------------------------------------------------
CREATE TABLE booking_cancellation (
  booking_cancellation_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  cancellation_reason_universal_uuid UUID DEFAULT NULL,
  cancellation_date DATE NOT NULL,
  refund_amount DECIMAL(18,2) DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_cancellation_uuid),
  CONSTRAINT uk_booking_cancellation_booking UNIQUE (booking_uuid)
);

CREATE INDEX idx_booking_cancellation_tenant ON booking_cancellation (tenant_uuid);
CREATE INDEX idx_booking_cancellation_workspace ON booking_cancellation (workspace_uuid);
CREATE INDEX idx_booking_cancellation_booking ON booking_cancellation (booking_uuid);
CREATE INDEX idx_booking_cancellation_is_deleted ON booking_cancellation (is_deleted);
CREATE INDEX idx_booking_cancellation_reason_univ ON booking_cancellation (cancellation_reason_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `booking_charge`
-- --------------------------------------------------------
CREATE TABLE booking_charge (
  booking_charge_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  charge_type_universal_uuid UUID DEFAULT NULL,
  charge_amount DECIMAL(18,2) NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_charge_uuid),
  CONSTRAINT uk_booking_charge_type_universal UNIQUE (booking_uuid, charge_type_universal_uuid)
);

CREATE INDEX idx_booking_charge_tenant ON booking_charge (tenant_uuid);
CREATE INDEX idx_booking_charge_workspace ON booking_charge (workspace_uuid);
CREATE INDEX idx_booking_charge_booking ON booking_charge (booking_uuid);
CREATE INDEX idx_booking_charge_is_deleted ON booking_charge (is_deleted);
CREATE INDEX idx_booking_charge_type_univ ON booking_charge (charge_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `booking_discount`
-- --------------------------------------------------------
CREATE TABLE booking_discount (
  booking_discount_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  discount_type_universal_uuid UUID DEFAULT NULL,
  discount_value DECIMAL(18,2) NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_discount_uuid),
  CONSTRAINT uk_booking_discount_type_universal UNIQUE (booking_uuid, discount_type_universal_uuid)
);

CREATE INDEX idx_booking_discount_tenant ON booking_discount (tenant_uuid);
CREATE INDEX idx_booking_discount_workspace ON booking_discount (workspace_uuid);
CREATE INDEX idx_booking_discount_booking ON booking_discount (booking_uuid);
CREATE INDEX idx_booking_discount_is_deleted ON booking_discount (is_deleted);
CREATE INDEX idx_booking_discount_type_univ ON booking_discount (discount_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `booking_parking`
-- --------------------------------------------------------
CREATE TABLE booking_parking (
  booking_parking_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  parking_uuid UUID NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_parking_uuid),
  CONSTRAINT uk_booking_parking_unique UNIQUE (parking_uuid)
);

CREATE INDEX idx_booking_parking_tenant ON booking_parking (tenant_uuid);
CREATE INDEX idx_booking_parking_workspace ON booking_parking (workspace_uuid);
CREATE INDEX idx_booking_parking_booking ON booking_parking (booking_uuid);
CREATE INDEX idx_booking_parking_item ON booking_parking (parking_uuid);
CREATE INDEX idx_booking_parking_is_deleted ON booking_parking (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `booking_payment`
-- --------------------------------------------------------
CREATE TABLE booking_payment (
  booking_payment_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  payment_status_universal_uuid UUID DEFAULT NULL,
  payment_mode_universal_uuid UUID DEFAULT NULL,
  receipt_number VARCHAR(100) NOT NULL,
  receipt_date DATE NOT NULL,
  amount DECIMAL(18,2) NOT NULL,
  reference_number VARCHAR(255) DEFAULT NULL,
  bank_name VARCHAR(255) DEFAULT NULL,
  cheque_date DATE DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_payment_uuid),
  CONSTRAINT uk_workspace_receipt_number UNIQUE (workspace_uuid, receipt_number)
);

CREATE INDEX idx_booking_payment_tenant ON booking_payment (tenant_uuid);
CREATE INDEX idx_booking_payment_workspace ON booking_payment (workspace_uuid);
CREATE INDEX idx_booking_payment_booking ON booking_payment (booking_uuid);
CREATE INDEX idx_booking_payment_receipt_date ON booking_payment (receipt_date);
CREATE INDEX idx_booking_payment_is_deleted ON booking_payment (is_deleted);
CREATE INDEX idx_booking_payment_calc ON booking_payment (booking_uuid, amount);

-- --------------------------------------------------------
-- Table structure for table `booking_refund`
-- --------------------------------------------------------
CREATE TABLE booking_refund (
  booking_refund_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  payment_mode_universal_uuid UUID DEFAULT NULL,
  refund_number VARCHAR(100) NOT NULL,
  refund_date DATE NOT NULL,
  refund_amount DECIMAL(18,2) NOT NULL,
  reference_number VARCHAR(255) DEFAULT NULL,
  bank_name VARCHAR(255) DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_refund_uuid),
  CONSTRAINT uk_workspace_refund_number UNIQUE (workspace_uuid, refund_number)
);

CREATE INDEX idx_booking_refund_tenant ON booking_refund (tenant_uuid);
CREATE INDEX idx_booking_refund_workspace ON booking_refund (workspace_uuid);
CREATE INDEX idx_booking_refund_booking ON booking_refund (booking_uuid);
CREATE INDEX idx_booking_refund_date ON booking_refund (refund_date);
CREATE INDEX idx_booking_refund_is_deleted ON booking_refund (is_deleted);
CREATE INDEX idx_booking_refund_mode_univ ON booking_refund (payment_mode_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `booking_storage`
-- --------------------------------------------------------
CREATE TABLE booking_storage (
  booking_storage_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  booking_uuid UUID NOT NULL,
  storage_uuid UUID NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (booking_storage_uuid),
  CONSTRAINT uk_booking_storage_unique UNIQUE (storage_uuid)
);

CREATE INDEX idx_booking_storage_tenant ON booking_storage (tenant_uuid);
CREATE INDEX idx_booking_storage_workspace ON booking_storage (workspace_uuid);
CREATE INDEX idx_booking_storage_booking ON booking_storage (booking_uuid);
CREATE INDEX idx_booking_storage_item ON booking_storage (storage_uuid);
CREATE INDEX idx_booking_storage_is_deleted ON booking_storage (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `calling_list`
-- --------------------------------------------------------
CREATE TABLE calling_list (
  calling_list_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  campaign_uuid UUID DEFAULT NULL,
  campaign_import_uuid UUID DEFAULT NULL,
  source_universal_uuid UUID DEFAULT NULL,
  calling_list_name VARCHAR(255) NOT NULL,
  assignment_type_universal_uuid UUID DEFAULT NULL,
  assigned_role_uuid UUID DEFAULT NULL,
  assigned_user_uuid UUID DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (calling_list_uuid)
);

CREATE INDEX idx_calling_list_tenant ON calling_list (tenant_uuid);
CREATE INDEX idx_calling_list_workspace ON calling_list (workspace_uuid);
CREATE INDEX idx_calling_list_campaign ON calling_list (campaign_uuid);
CREATE INDEX idx_calling_list_import ON calling_list (campaign_import_uuid);
CREATE INDEX idx_calling_list_assigned_user ON calling_list (assigned_user_uuid);
CREATE INDEX idx_calling_list_is_deleted ON calling_list (is_deleted);
CREATE INDEX idx_calling_list_source_univ ON calling_list (source_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `calling_list_record`
-- --------------------------------------------------------
CREATE TABLE calling_list_record (
  calling_list_record_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  calling_list_uuid UUID NOT NULL,
  contact_uuid UUID NOT NULL,
  lead_uuid UUID DEFAULT NULL,
  opportunity_uuid UUID DEFAULT NULL,
  assigned_user_uuid UUID DEFAULT NULL,
  phone_number VARCHAR(20) NOT NULL,
  next_call_date DATE DEFAULT NULL,
  last_call_at TIMESTAMPTZ DEFAULT NULL,
  last_call_status_uuid UUID DEFAULT NULL,
  call_count INT NOT NULL DEFAULT 0,
  is_closed BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (calling_list_record_uuid)
);

CREATE INDEX idx_calling_list_rec_tenant ON calling_list_record (tenant_uuid);
CREATE INDEX idx_calling_list_rec_workspace ON calling_list_record (workspace_uuid);
CREATE INDEX idx_calling_list_rec_list ON calling_list_record (calling_list_uuid);
CREATE INDEX idx_calling_list_rec_contact ON calling_list_record (contact_uuid);
CREATE INDEX idx_calling_list_rec_lead ON calling_list_record (lead_uuid);
CREATE INDEX idx_calling_list_rec_opp ON calling_list_record (opportunity_uuid);
CREATE INDEX idx_calling_list_rec_user ON calling_list_record (assigned_user_uuid);
CREATE INDEX idx_calling_list_rec_phone ON calling_list_record (phone_number);
CREATE INDEX idx_calling_list_rec_next_call ON calling_list_record (next_call_date);
CREATE INDEX idx_calling_list_rec_last_call ON calling_list_record (last_call_at);
CREATE INDEX idx_calling_list_rec_status ON calling_list_record (last_call_status_uuid);
CREATE INDEX idx_calling_list_rec_closed ON calling_list_record (is_closed);
CREATE INDEX idx_calling_list_rec_deleted ON calling_list_record (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `call_log`
-- --------------------------------------------------------
CREATE TABLE call_log (
  call_log_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID DEFAULT NULL,
  user_uuid UUID NOT NULL,
  calling_list_record_uuid UUID DEFAULT NULL,
  activity_uuid UUID DEFAULT NULL,
  contact_uuid UUID DEFAULT NULL,
  lead_uuid UUID DEFAULT NULL,
  opportunity_uuid UUID DEFAULT NULL,
  phone_number VARCHAR(20) NOT NULL,
  call_direction_universal_uuid UUID DEFAULT NULL,
  call_status_universal_uuid UUID DEFAULT NULL,
  call_start_datetime TIMESTAMPTZ NOT NULL,
  call_date DATE NOT NULL,
  call_end_datetime TIMESTAMPTZ DEFAULT NULL,
  duration_seconds INT NOT NULL DEFAULT 0,
  is_known_contact BOOLEAN NOT NULL DEFAULT FALSE,
  is_from_calling_list BOOLEAN NOT NULL DEFAULT FALSE,
  device_call_log_id VARCHAR(255) DEFAULT NULL,
  sim_slot VARCHAR(50) DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (call_log_uuid),
  CONSTRAINT uk_user_device_call_log UNIQUE (user_uuid, device_call_log_id)
);

CREATE INDEX idx_call_log_tenant ON call_log (tenant_uuid);
CREATE INDEX idx_call_log_workspace ON call_log (workspace_uuid);
CREATE INDEX idx_call_log_user ON call_log (user_uuid);
CREATE INDEX idx_call_log_record ON call_log (calling_list_record_uuid);
CREATE INDEX idx_call_log_activity ON call_log (activity_uuid);
CREATE INDEX idx_call_log_contact ON call_log (contact_uuid);
CREATE INDEX idx_call_log_lead ON call_log (lead_uuid);
CREATE INDEX idx_call_log_opp ON call_log (opportunity_uuid);
CREATE INDEX idx_call_log_phone ON call_log (phone_number);
CREATE INDEX idx_call_log_start_time ON call_log (call_start_datetime);
CREATE INDEX idx_call_log_date ON call_log (call_date);
CREATE INDEX idx_call_log_known ON call_log (is_known_contact);
CREATE INDEX idx_call_log_from_list ON call_log (is_from_calling_list);
CREATE INDEX idx_call_log_deleted ON call_log (is_deleted);
CREATE INDEX idx_call_log_same_day_check ON call_log (tenant_uuid, workspace_uuid, phone_number, call_date);

-- --------------------------------------------------------
-- Table structure for table `campaign`
-- --------------------------------------------------------
CREATE TABLE campaign (
  campaign_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  campaign_type_universal_uuid UUID DEFAULT NULL,
  status_uuid UUID DEFAULT NULL,
  campaign_name VARCHAR(255) NOT NULL,
  start_date DATE DEFAULT NULL,
  end_date DATE DEFAULT NULL,
  estimated_budget DECIMAL(18,2) DEFAULT NULL,
  actual_budget DECIMAL(18,2) DEFAULT NULL,
  campaign_settings JSONB DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (campaign_uuid),
  CONSTRAINT uk_workspace_campaign_name UNIQUE (tenant_uuid, workspace_uuid, campaign_name)
);

CREATE INDEX idx_campaign_tenant_uuid ON campaign (tenant_uuid);
CREATE INDEX idx_campaign_workspace_uuid ON campaign (workspace_uuid);
CREATE INDEX idx_campaign_status_uuid ON campaign (status_uuid);
CREATE INDEX idx_campaign_is_deleted ON campaign (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `campaign_import`
-- --------------------------------------------------------
CREATE TABLE campaign_import (
  campaign_import_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  campaign_uuid UUID NOT NULL,
  file_name VARCHAR(255) DEFAULT NULL,
  total_records INT DEFAULT NULL,
  successful_records INT DEFAULT NULL,
  failed_records INT DEFAULT NULL,
  import_summary JSONB DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (campaign_import_uuid)
);

CREATE INDEX idx_campaign_imp_tenant ON campaign_import (tenant_uuid);
CREATE INDEX idx_campaign_imp_workspace ON campaign_import (workspace_uuid);
CREATE INDEX idx_campaign_imp_campaign ON campaign_import (campaign_uuid);
CREATE INDEX idx_campaign_imp_deleted ON campaign_import (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `campaign_member`
-- --------------------------------------------------------
CREATE TABLE campaign_member (
  campaign_member_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  campaign_uuid UUID NOT NULL,
  contact_uuid UUID NOT NULL,
  lead_uuid UUID DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (campaign_member_uuid),
  CONSTRAINT uk_campaign_contact UNIQUE (campaign_uuid, contact_uuid)
);

CREATE INDEX idx_campaign_mem_tenant ON campaign_member (tenant_uuid);
CREATE INDEX idx_campaign_mem_workspace ON campaign_member (workspace_uuid);
CREATE INDEX idx_campaign_mem_campaign ON campaign_member (campaign_uuid);
CREATE INDEX idx_campaign_mem_contact ON campaign_member (contact_uuid);
CREATE INDEX idx_campaign_mem_lead ON campaign_member (lead_uuid);
CREATE INDEX idx_campaign_mem_deleted ON campaign_member (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `commission`
-- --------------------------------------------------------
CREATE TABLE commission (
  commission_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  commission_type VARCHAR(20) NOT NULL,
  commission_percentage DECIMAL(18,2) DEFAULT NULL,
  commission_amount DECIMAL(18,2) DEFAULT NULL,
  status_uuid UUID DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (commission_uuid),
  CONSTRAINT uk_entity_receiving_user UNIQUE (module_uuid, entity_uuid, user_uuid)
);

CREATE INDEX idx_commission_tenant ON commission (tenant_uuid);
CREATE INDEX idx_commission_workspace ON commission (workspace_uuid);
CREATE INDEX idx_commission_module_entity ON commission (module_uuid, entity_uuid);
CREATE INDEX idx_commission_user ON commission (user_uuid);
CREATE INDEX idx_commission_status ON commission (status_uuid);
CREATE INDEX idx_commission_deleted ON commission (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `commission_payment`
-- --------------------------------------------------------
CREATE TABLE commission_payment (
  commission_payment_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  commission_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  payment_mode_universal_uuid UUID DEFAULT NULL,
  payment_status_universal_uuid UUID DEFAULT NULL,
  payment_number VARCHAR(100) NOT NULL,
  payment_date DATE NOT NULL,
  amount DECIMAL(18,2) NOT NULL,
  reference_number VARCHAR(255) DEFAULT NULL,
  bank_name VARCHAR(255) DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (commission_payment_uuid),
  CONSTRAINT uk_workspace_payment_number UNIQUE (workspace_uuid, payment_number)
);

CREATE INDEX idx_comm_pay_tenant ON commission_payment (tenant_uuid);
CREATE INDEX idx_comm_pay_workspace ON commission_payment (workspace_uuid);
CREATE INDEX idx_comm_pay_commission ON commission_payment (commission_uuid);
CREATE INDEX idx_comm_pay_user ON commission_payment (user_uuid);
CREATE INDEX idx_comm_pay_date ON commission_payment (payment_date);
CREATE INDEX idx_comm_pay_deleted ON commission_payment (is_deleted);
CREATE INDEX idx_comm_pay_mode_univ ON commission_payment (payment_mode_universal_uuid);
CREATE INDEX idx_comm_pay_status_univ ON commission_payment (payment_status_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `contact`
-- --------------------------------------------------------
CREATE TABLE contact (
  contact_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  owner_user_uuid UUID DEFAULT NULL,
  full_name VARCHAR(255) NOT NULL,
  mobile_numbers JSONB DEFAULT NULL,
  email_addresses JSONB DEFAULT NULL,
  addresses JSONB DEFAULT NULL,
  relationship_details JSONB DEFAULT NULL,
  contact_search_text TEXT DEFAULT NULL,
  company_name VARCHAR(255) DEFAULT NULL,
  occupation VARCHAR(255) DEFAULT NULL,
  birth_date DATE DEFAULT NULL,
  anniversary_date DATE DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (contact_uuid)
);

CREATE INDEX idx_contact_tenant_uuid ON contact (tenant_uuid);
CREATE INDEX idx_contact_workspace_uuid ON contact (workspace_uuid);
CREATE INDEX idx_contact_owner_user_uuid ON contact (owner_user_uuid);
CREATE INDEX idx_contact_full_name ON contact (full_name);
CREATE INDEX idx_contact_is_deleted ON contact (is_deleted);
CREATE INDEX idx_contact_mobile_numbers ON contact USING gin (mobile_numbers);

-- --------------------------------------------------------
-- Table structure for table `document`
-- --------------------------------------------------------
CREATE TABLE document (
  document_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  document_type_universal_uuid UUID DEFAULT NULL,
  document_title VARCHAR(255) DEFAULT NULL,
  original_file_name VARCHAR(255) NOT NULL,
  stored_file_name VARCHAR(255) NOT NULL,
  file_extension VARCHAR(20) DEFAULT NULL,
  mime_type VARCHAR(100) DEFAULT NULL,
  file_size BIGINT DEFAULT NULL,
  file_path VARCHAR(1000) NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (document_uuid)
);

CREATE INDEX idx_document_tenant_uuid ON document (tenant_uuid);
CREATE INDEX idx_document_module_uuid ON document (module_uuid);
CREATE INDEX idx_document_entity_uuid ON document (entity_uuid);
CREATE INDEX idx_document_is_deleted ON document (is_deleted);
CREATE INDEX idx_document_module_entity ON document (module_uuid, entity_uuid);
CREATE INDEX idx_document_module_entity_created_at ON document (module_uuid, entity_uuid, created_at);
CREATE INDEX idx_document_type_universal_uuid ON document (document_type_universal_uuid);
CREATE INDEX idx_document_combo ON document (module_uuid, entity_uuid, document_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `lead`
-- --------------------------------------------------------
CREATE TABLE lead (
  lead_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  contact_uuid UUID NOT NULL,
  source_universal_uuid UUID DEFAULT NULL,
  campaign_uuid UUID DEFAULT NULL,
  property_uuid UUID DEFAULT NULL,
  assigned_user_uuid UUID DEFAULT NULL,
  visibility_user_uuid UUID DEFAULT NULL,
  enquiry_date TIMESTAMPTZ DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (lead_uuid)
);

CREATE INDEX idx_lead_tenant_uuid ON lead (tenant_uuid);
CREATE INDEX idx_lead_workspace_uuid ON lead (workspace_uuid);
CREATE INDEX idx_lead_contact_uuid ON lead (contact_uuid);
CREATE INDEX idx_lead_campaign_uuid ON lead (campaign_uuid);
CREATE INDEX idx_lead_assigned_user_uuid ON lead (assigned_user_uuid);
CREATE INDEX idx_lead_is_deleted ON lead (is_deleted);
CREATE INDEX idx_lead_visibility_user_uuid ON lead (visibility_user_uuid);
CREATE INDEX idx_lead_source_universal_uuid ON lead (source_universal_uuid);
CREATE INDEX idx_lead_property_uuid ON lead (property_uuid);

-- --------------------------------------------------------
-- Table structure for table `media`
-- --------------------------------------------------------
CREATE TABLE media (
  media_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  module_uuid UUID DEFAULT NULL,
  entity_uuid UUID DEFAULT NULL,
  property_block_uuid UUID DEFAULT NULL,
  unit_uuid UUID DEFAULT NULL,
  property_uuid UUID DEFAULT NULL,
  media_type_uuid UUID NOT NULL,
  title VARCHAR(255) DEFAULT NULL,
  description TEXT DEFAULT NULL,
  file_name VARCHAR(255) DEFAULT NULL,
  original_file_name VARCHAR(255) DEFAULT NULL,
  file_extension VARCHAR(20) DEFAULT NULL,
  mime_type VARCHAR(100) DEFAULT NULL,
  file_size BIGINT DEFAULT NULL,
  file_url TEXT DEFAULT NULL,
  thumbnail_url TEXT DEFAULT NULL,
  youtube_url VARCHAR(1000) DEFAULT NULL,
  external_url VARCHAR(1000) DEFAULT NULL,
  embed_url VARCHAR(1000) DEFAULT NULL,
  display_order INT NOT NULL DEFAULT 0,
  is_featured BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  tags_json JSONB DEFAULT NULL,
  PRIMARY KEY (media_uuid)
);

CREATE INDEX idx_media_tenant ON media (tenant_uuid);
CREATE INDEX idx_media_workspace ON media (workspace_uuid);
CREATE INDEX idx_media_module_entity ON media (module_uuid, entity_uuid);
CREATE INDEX idx_media_property_block ON media (property_block_uuid);
CREATE INDEX idx_media_unit ON media (unit_uuid);
CREATE INDEX idx_media_type ON media (media_type_uuid);
CREATE INDEX idx_media_featured ON media (is_featured);
CREATE INDEX idx_media_active ON media (is_active);
CREATE INDEX idx_media_property_uuid ON media (property_uuid);

-- --------------------------------------------------------
-- Table structure for table `message_log`
-- --------------------------------------------------------
CREATE TABLE message_log (
  message_log_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID DEFAULT NULL,
  module_uuid UUID DEFAULT NULL,
  entity_uuid UUID DEFAULT NULL,
  message_channel_universal_uuid UUID DEFAULT NULL,
  recipient VARCHAR(255) NOT NULL,
  subject VARCHAR(500) DEFAULT NULL,
  message_body TEXT DEFAULT NULL,
  external_reference VARCHAR(255) DEFAULT NULL,
  status_uuid UUID DEFAULT NULL,
  sent_at TIMESTAMPTZ DEFAULT NULL,
  delivered_at TIMESTAMPTZ DEFAULT NULL,
  failed_at TIMESTAMPTZ DEFAULT NULL,
  error_message VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  PRIMARY KEY (message_log_uuid)
);

CREATE INDEX idx_message_log_tenant_uuid ON message_log (tenant_uuid);
CREATE INDEX idx_message_log_module_entity ON message_log (module_uuid, entity_uuid);
CREATE INDEX idx_message_log_sent_at ON message_log (sent_at);
CREATE INDEX idx_message_log_workspace_uuid ON message_log (workspace_uuid);
CREATE INDEX idx_message_log_channel_universal_uuid ON message_log (message_channel_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `message_queue`
-- --------------------------------------------------------
CREATE TABLE message_queue (
  message_queue_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID DEFAULT NULL,
  module_uuid UUID DEFAULT NULL,
  entity_uuid UUID DEFAULT NULL,
  message_channel_universal_uuid UUID DEFAULT NULL,
  message_template_uuid UUID DEFAULT NULL,
  recipient VARCHAR(255) DEFAULT NULL,
  subject VARCHAR(500) DEFAULT NULL,
  message_body TEXT DEFAULT NULL,
  payload_json JSONB DEFAULT NULL,
  retry_count INT NOT NULL DEFAULT 0,
  sent_at TIMESTAMPTZ DEFAULT NULL,
  failed_at TIMESTAMPTZ DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  PRIMARY KEY (message_queue_uuid)
);

CREATE INDEX idx_message_queue_sent_at ON message_queue (sent_at);
CREATE INDEX idx_message_queue_module_entity ON message_queue (module_uuid, entity_uuid);
CREATE INDEX idx_message_queue_sent_at_created_at ON message_queue (sent_at, created_at);
CREATE INDEX idx_message_queue_module_entity_created_at ON message_queue (module_uuid, entity_uuid, created_at);
CREATE INDEX idx_message_queue_workspace_uuid ON message_queue (workspace_uuid);
CREATE INDEX idx_message_queue_channel_universal_uuid ON message_queue (message_channel_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `message_template`
-- --------------------------------------------------------
CREATE TABLE message_template (
  message_template_uuid UUID NOT NULL,
  tenant_uuid UUID DEFAULT NULL,
  workspace_uuid UUID DEFAULT NULL,
  message_channel_universal_uuid UUID DEFAULT NULL,
  template_name VARCHAR(255) NOT NULL,
  template_subject VARCHAR(500) DEFAULT NULL,
  template_body TEXT NOT NULL,
  variable_list_json JSONB DEFAULT NULL,
  is_system BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (message_template_uuid),
  CONSTRAINT uk_tenant_workspace_template UNIQUE (tenant_uuid, workspace_uuid, message_channel_universal_uuid, template_name)
);

CREATE INDEX idx_msg_temp_tenant_uuid ON message_template (tenant_uuid);
CREATE INDEX idx_msg_temp_workspace_uuid ON message_template (workspace_uuid);
CREATE INDEX idx_msg_temp_channel_univ ON message_template (message_channel_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `mst_currency`
-- --------------------------------------------------------
CREATE TABLE mst_currency (
  currency_uuid UUID NOT NULL,
  currency_code VARCHAR(10) NOT NULL,
  currency_name VARCHAR(100) NOT NULL,
  currency_symbol VARCHAR(10) DEFAULT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  is_system BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (currency_uuid),
  CONSTRAINT uk_currency_code UNIQUE (currency_code)
);

-- --------------------------------------------------------
-- Table structure for table `mst_master_type`
-- --------------------------------------------------------
CREATE TABLE mst_master_type (
  master_type_uuid UUID NOT NULL,
  master_type_name VARCHAR(150) NOT NULL,
  display_order INT NOT NULL DEFAULT 0,
  is_system BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (master_type_uuid),
  CONSTRAINT uk_master_type_name UNIQUE (master_type_name)
);

-- --------------------------------------------------------
-- Table structure for table `mst_module`
-- --------------------------------------------------------
CREATE TABLE mst_module (
  module_uuid UUID NOT NULL,
  module_name VARCHAR(100) NOT NULL,
  is_system BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (module_uuid),
  CONSTRAINT uk_module_name UNIQUE (module_name)
);

-- --------------------------------------------------------
-- Table structure for table `mst_permission`
-- --------------------------------------------------------
CREATE TABLE mst_permission (
  permission_uuid UUID NOT NULL,
  permission_code VARCHAR(100) NOT NULL,
  permission_name VARCHAR(255) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (permission_uuid),
  CONSTRAINT uk_permission_code UNIQUE (permission_code)
);

-- --------------------------------------------------------
-- Table structure for table `mst_plan`
-- --------------------------------------------------------
CREATE TABLE mst_plan (
  plan_uuid UUID NOT NULL,
  plan_type_universal_uuid UUID DEFAULT NULL,
  plan_name VARCHAR(100) NOT NULL,
  max_users INT NOT NULL DEFAULT 0,
  default_storage_bytes BIGINT NOT NULL DEFAULT 0,
  default_sms_credits INT NOT NULL DEFAULT 0,
  default_whatsapp_credits INT NOT NULL DEFAULT 0,
  default_email_credits INT NOT NULL DEFAULT 0,
  sort_order INT NOT NULL DEFAULT 0,
  is_system BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (plan_uuid),
  CONSTRAINT uk_plan_type_plan_name UNIQUE (plan_type_universal_uuid, plan_name)
);

CREATE INDEX idx_plan_name ON mst_plan (plan_name);
CREATE INDEX idx_plan_type_universal_uuid ON mst_plan (plan_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `mst_plan_price`
-- --------------------------------------------------------
CREATE TABLE mst_plan_price (
  plan_price_uuid UUID NOT NULL,
  plan_uuid UUID NOT NULL,
  billing_cycle_universal_uuid UUID DEFAULT NULL,
  price_amount DECIMAL(18,2) NOT NULL,
  effective_from_date DATE NOT NULL,
  effective_to_date DATE DEFAULT NULL,
  is_force_migration BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (plan_price_uuid)
);

CREATE INDEX idx_plan_price_plan_uuid ON mst_plan_price (plan_uuid);
CREATE INDEX idx_plan_price_eff_from ON mst_plan_price (effective_from_date);
CREATE INDEX idx_plan_price_cycle_univ ON mst_plan_price (billing_cycle_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `mst_status`
-- --------------------------------------------------------
CREATE TABLE mst_status (
  status_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  status_name VARCHAR(100) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  color_code VARCHAR(20) DEFAULT NULL,
  is_system BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (status_uuid),
  CONSTRAINT uk_module_status UNIQUE (module_uuid, status_name)
);

CREATE INDEX idx_mst_status_module_uuid ON mst_status (module_uuid);

-- --------------------------------------------------------
-- Table structure for table `mst_universal`
-- --------------------------------------------------------
CREATE TABLE mst_universal (
  universal_uuid UUID NOT NULL,
  tenant_uuid UUID DEFAULT NULL,
  master_type_uuid UUID NOT NULL,
  parent_universal_uuid UUID DEFAULT NULL,
  value_name VARCHAR(255) NOT NULL,
  display_order INT NOT NULL DEFAULT 0,
  is_system BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (universal_uuid),
  CONSTRAINT uk_master_value UNIQUE (tenant_uuid, master_type_uuid, parent_universal_uuid, value_name)
);

CREATE INDEX idx_mst_universal_tenant_uuid ON mst_universal (tenant_uuid);
CREATE INDEX idx_mst_universal_master_type_uuid ON mst_universal (master_type_uuid);
CREATE INDEX idx_mst_universal_parent_universal_uuid ON mst_universal (parent_universal_uuid);
CREATE INDEX idx_mst_universal_is_deleted ON mst_universal (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `note`
-- --------------------------------------------------------
CREATE TABLE note (
  note_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  note_text TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (note_uuid)
);

CREATE INDEX idx_note_tenant_uuid ON note (tenant_uuid);
CREATE INDEX idx_note_module_uuid ON note (module_uuid);
CREATE INDEX idx_note_entity_uuid ON note (entity_uuid);
CREATE INDEX idx_note_is_deleted ON note (is_deleted);
CREATE INDEX idx_note_module_entity ON note (module_uuid, entity_uuid);
CREATE INDEX idx_note_module_entity_created_at ON note (module_uuid, entity_uuid, created_at);

-- --------------------------------------------------------
-- Table structure for table `opportunity`
-- --------------------------------------------------------
CREATE TABLE opportunity (
  opportunity_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  contact_uuid UUID NOT NULL,
  lead_uuid UUID NOT NULL,
  property_uuid UUID DEFAULT NULL,
  assigned_user_uuid UUID DEFAULT NULL,
  opportunity_stage_universal_uuid UUID DEFAULT NULL,
  expected_budget DECIMAL(18,2) DEFAULT NULL,
  expected_closing_date DATE DEFAULT NULL,
  next_followup_datetime TIMESTAMPTZ DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (opportunity_uuid)
);

CREATE INDEX idx_opportunity_tenant_uuid ON opportunity (tenant_uuid);
CREATE INDEX idx_opportunity_workspace_uuid ON opportunity (workspace_uuid);
CREATE INDEX idx_opportunity_contact_uuid ON opportunity (contact_uuid);
CREATE INDEX idx_opportunity_lead_uuid ON opportunity (lead_uuid);
CREATE INDEX idx_opportunity_assigned_user_uuid ON opportunity (assigned_user_uuid);
CREATE INDEX idx_opportunity_next_followup_datetime ON opportunity (next_followup_datetime);
CREATE INDEX idx_opportunity_is_deleted ON opportunity (is_deleted);
CREATE INDEX idx_opportunity_stage_universal_uuid ON opportunity (opportunity_stage_universal_uuid);
CREATE INDEX idx_opportunity_property_uuid ON opportunity (property_uuid);

-- --------------------------------------------------------
-- Table structure for table `parking`
-- --------------------------------------------------------
CREATE TABLE parking (
  parking_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  property_block_uuid UUID DEFAULT NULL,
  parking_type_universal_uuid UUID DEFAULT NULL,
  parking_number VARCHAR(100) NOT NULL,
  is_booked BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (parking_uuid),
  CONSTRAINT uk_workspace_parking_number UNIQUE (workspace_uuid, parking_number)
);

CREATE INDEX idx_parking_property_block_uuid ON parking (property_block_uuid);
CREATE INDEX idx_parking_is_deleted ON parking (is_deleted);
CREATE INDEX idx_parking_type_universal_uuid ON parking (parking_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `property`
-- --------------------------------------------------------
CREATE TABLE property (
  property_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  owner_user_uuid UUID DEFAULT NULL,
  property_title VARCHAR(255) NOT NULL,
  public_slug VARCHAR(255) DEFAULT NULL,
  public_url_token VARCHAR(100) DEFAULT NULL,
  is_public_listing BOOLEAN NOT NULL DEFAULT FALSE,
  unit_type_universal_uuid UUID DEFAULT NULL,
  property_status_universal_uuid UUID DEFAULT NULL,
  source_universal_uuid UUID DEFAULT NULL,
  listing_type VARCHAR(20) NOT NULL,
  expected_price DECIMAL(18,2) DEFAULT NULL,
  expected_rent DECIMAL(18,2) DEFAULT NULL,
  carpet_area DECIMAL(18,2) DEFAULT NULL,
  builtup_area DECIMAL(18,2) DEFAULT NULL,
  saleable_area DECIMAL(18,2) DEFAULT NULL,
  area_unit_universal_uuid UUID DEFAULT NULL,
  floor_number INT DEFAULT NULL,
  total_floors INT DEFAULT NULL,
  facing_details JSONB DEFAULT NULL,
  amenity_details JSONB DEFAULT NULL,
  address_details JSONB DEFAULT NULL,
  contact_person_details JSONB DEFAULT NULL,
  available_from DATE DEFAULT NULL,
  remarks VARCHAR(2000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (property_uuid)
);

CREATE INDEX idx_property_tenant_uuid ON property (tenant_uuid);
CREATE INDEX idx_property_workspace_uuid ON property (workspace_uuid);
CREATE INDEX idx_property_owner_user_uuid ON property (owner_user_uuid);
CREATE INDEX idx_property_unit_type_universal_uuid ON property (unit_type_universal_uuid);
CREATE INDEX idx_property_status_universal_uuid ON property (property_status_universal_uuid);
CREATE INDEX idx_property_source_universal_uuid ON property (source_universal_uuid);
CREATE INDEX idx_property_listing_type ON property (listing_type);
CREATE INDEX idx_property_is_deleted ON property (is_deleted);
CREATE INDEX idx_property_public_slug ON property (public_slug);
CREATE INDEX idx_property_public_url_token ON property (public_url_token);
CREATE INDEX idx_property_is_public_listing ON property (is_public_listing);

-- --------------------------------------------------------
-- Table structure for table `property_block`
-- --------------------------------------------------------
CREATE TABLE property_block (
  property_block_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  property_block_type_universal_uuid UUID DEFAULT NULL,
  property_block_name VARCHAR(255) NOT NULL,
  property_block_code VARCHAR(50) DEFAULT NULL,
  rera_number VARCHAR(100) DEFAULT NULL,
  number_of_floors INT DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (property_block_uuid),
  CONSTRAINT uk_workspace_property_block_name UNIQUE (workspace_uuid, property_block_name)
);

CREATE INDEX idx_property_block_tenant ON property_block (tenant_uuid);
CREATE INDEX idx_property_block_workspace ON property_block (workspace_uuid);
CREATE INDEX idx_property_block_name ON property_block (property_block_name);
CREATE INDEX idx_property_block_is_deleted ON property_block (is_deleted);
CREATE INDEX idx_property_block_type_univ ON property_block (property_block_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `property_block_stage_progress`
-- --------------------------------------------------------
CREATE TABLE property_block_stage_progress (
  property_block_stage_progress_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  property_block_uuid UUID NOT NULL,
  construction_stage_universal_uuid UUID DEFAULT NULL,
  completion_date DATE DEFAULT NULL,
  architect_certificate_document_uuid UUID DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (property_block_stage_progress_uuid),
  CONSTRAINT uk_property_block_stage UNIQUE (property_block_uuid, construction_stage_universal_uuid)
);

CREATE INDEX idx_pb_stage_tenant ON property_block_stage_progress (tenant_uuid);
CREATE INDEX idx_pb_stage_workspace ON property_block_stage_progress (workspace_uuid);
CREATE INDEX idx_pb_stage_block ON property_block_stage_progress (property_block_uuid);
CREATE INDEX idx_pb_stage_is_deleted ON property_block_stage_progress (is_deleted);
CREATE INDEX idx_pb_stage_construction_univ ON property_block_stage_progress (construction_stage_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `property_transaction`
-- --------------------------------------------------------
CREATE TABLE property_transaction (
  property_transaction_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  property_uuid UUID NOT NULL,
  opportunity_uuid UUID DEFAULT NULL,
  quotation_uuid UUID DEFAULT NULL,
  contact_uuid UUID NOT NULL,
  status_uuid UUID DEFAULT NULL,
  transaction_number VARCHAR(100) NOT NULL,
  transaction_type VARCHAR(20) NOT NULL,
  transaction_date DATE NOT NULL,
  deal_amount DECIMAL(18,2) DEFAULT NULL,
  monthly_rent DECIMAL(18,2) DEFAULT NULL,
  deposit_amount DECIMAL(18,2) DEFAULT NULL,
  brokerage_amount DECIMAL(18,2) DEFAULT NULL,
  brokerage_received_amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  brokerage_balance_amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  agreement_date DATE DEFAULT NULL,
  possession_date DATE DEFAULT NULL,
  remarks VARCHAR(2000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (property_transaction_uuid),
  CONSTRAINT uk_workspace_property_transaction_number UNIQUE (workspace_uuid, transaction_number)
);

CREATE INDEX idx_property_transaction_tenant ON property_transaction (tenant_uuid);
CREATE INDEX idx_property_transaction_workspace ON property_transaction (workspace_uuid);
CREATE INDEX idx_property_transaction_property ON property_transaction (property_uuid);
CREATE INDEX idx_property_transaction_opportunity ON property_transaction (opportunity_uuid);
CREATE INDEX idx_property_transaction_quotation ON property_transaction (quotation_uuid);
CREATE INDEX idx_property_transaction_contact ON property_transaction (contact_uuid);
CREATE INDEX idx_property_transaction_status ON property_transaction (status_uuid);
CREATE INDEX idx_property_transaction_type ON property_transaction (transaction_type);
CREATE INDEX idx_property_transaction_date ON property_transaction (transaction_date);
CREATE INDEX idx_property_transaction_deleted ON property_transaction (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `quotation`
-- --------------------------------------------------------
CREATE TABLE quotation (
  quotation_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  opportunity_uuid UUID NOT NULL,
  contact_uuid UUID NOT NULL,
  status_uuid UUID DEFAULT NULL,
  quotation_number VARCHAR(100) NOT NULL,
  quotation_date DATE NOT NULL,
  valid_until DATE DEFAULT NULL,
  unit_uuid UUID DEFAULT NULL,
  property_uuid UUID DEFAULT NULL,
  total_amount DECIMAL(18,2) DEFAULT NULL,
  discount_amount DECIMAL(18,2) DEFAULT NULL,
  final_amount DECIMAL(18,2) DEFAULT NULL,
  quotation_snapshot TEXT DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (quotation_uuid),
  CONSTRAINT uk_workspace_quotation_number UNIQUE (workspace_uuid, quotation_number),
  CONSTRAINT uk_opportunity_quotation UNIQUE (opportunity_uuid)
);

CREATE INDEX idx_quotation_tenant_uuid ON quotation (tenant_uuid);
CREATE INDEX idx_quotation_workspace_uuid ON quotation (workspace_uuid);
CREATE INDEX idx_quotation_opportunity_uuid ON quotation (opportunity_uuid);
CREATE INDEX idx_quotation_contact_uuid ON quotation (contact_uuid);
CREATE INDEX idx_quotation_status_uuid ON quotation (status_uuid);
CREATE INDEX idx_quotation_unit_uuid ON quotation (unit_uuid);
CREATE INDEX idx_quotation_date ON quotation (quotation_date);
CREATE INDEX idx_quotation_is_deleted ON quotation (is_deleted);
CREATE INDEX idx_quotation_property_uuid ON quotation (property_uuid);

-- --------------------------------------------------------
-- Table structure for table `quotation_charge`
-- --------------------------------------------------------
CREATE TABLE quotation_charge (
  quotation_charge_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  quotation_uuid UUID NOT NULL,
  charge_type_universal_uuid UUID NOT NULL,
  charge_amount DECIMAL(18,2) NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (quotation_charge_uuid),
  CONSTRAINT uk_quotation_charge_type UNIQUE (quotation_uuid, charge_type_universal_uuid)
);

CREATE INDEX idx_quotation_charge_tenant ON quotation_charge (tenant_uuid);
CREATE INDEX idx_quotation_charge_workspace ON quotation_charge (workspace_uuid);
CREATE INDEX idx_quotation_charge_uuid ON quotation_charge (quotation_uuid);
CREATE INDEX idx_quotation_charge_type_univ ON quotation_charge (charge_type_universal_uuid);
CREATE INDEX idx_quotation_charge_is_deleted ON quotation_charge (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `quotation_discount`
-- --------------------------------------------------------
CREATE TABLE quotation_discount (
  quotation_discount_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  quotation_uuid UUID NOT NULL,
  discount_type_universal_uuid UUID NOT NULL,
  discount_value DECIMAL(18,2) NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (quotation_discount_uuid),
  CONSTRAINT uk_quotation_discount_type UNIQUE (quotation_uuid, discount_type_universal_uuid)
);

CREATE INDEX idx_quotation_discount_tenant ON quotation_discount (tenant_uuid);
CREATE INDEX idx_quotation_discount_workspace ON quotation_discount (workspace_uuid);
CREATE INDEX idx_quotation_discount_uuid ON quotation_discount (quotation_uuid);
CREATE INDEX idx_quotation_discount_type_univ ON quotation_discount (discount_type_universal_uuid);
CREATE INDEX idx_quotation_discount_is_deleted ON quotation_discount (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `role`
-- --------------------------------------------------------
CREATE TABLE role (
  role_uuid UUID NOT NULL,
  tenant_uuid UUID DEFAULT NULL,
  role_name VARCHAR(100) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  is_system BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (role_uuid),
  CONSTRAINT uk_role_name_system UNIQUE (tenant_uuid, role_name)
);

CREATE INDEX idx_role_tenant_uuid ON role (tenant_uuid);
CREATE INDEX idx_role_is_deleted ON role (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `role_permission`
-- --------------------------------------------------------
CREATE TABLE role_permission (
  role_permission_uuid UUID NOT NULL,
  tenant_uuid UUID DEFAULT NULL,
  role_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  permission_uuid UUID NOT NULL,
  is_granted BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (role_permission_uuid),
  CONSTRAINT uk_role_permission UNIQUE (tenant_uuid, role_uuid, module_uuid, permission_uuid),
  CONSTRAINT fk_role_permission_module FOREIGN KEY (module_uuid) REFERENCES mst_module (module_uuid)
);

CREATE INDEX idx_role_permission_tenant ON role_permission (tenant_uuid);
CREATE INDEX idx_role_permission_role ON role_permission (role_uuid);
CREATE INDEX idx_role_permission_module ON role_permission (module_uuid);
CREATE INDEX idx_role_permission_item ON role_permission (permission_uuid);
CREATE INDEX idx_role_permission_is_deleted ON role_permission (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `saved_filter`
-- --------------------------------------------------------
CREATE TABLE saved_filter (
  saved_filter_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID DEFAULT NULL,
  module_uuid UUID NOT NULL,
  user_uuid UUID DEFAULT NULL,
  filter_name VARCHAR(255) NOT NULL,
  filter_json JSONB NOT NULL,
  is_shared BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (saved_filter_uuid),
  CONSTRAINT uk_user_module_filter UNIQUE (tenant_uuid, workspace_uuid, module_uuid, user_uuid, filter_name)
);

CREATE INDEX idx_saved_filter_tenant_uuid ON saved_filter (tenant_uuid);
CREATE INDEX idx_saved_filter_workspace_uuid ON saved_filter (workspace_uuid);
CREATE INDEX idx_saved_filter_module_uuid ON saved_filter (module_uuid);
CREATE INDEX idx_saved_filter_user_uuid ON saved_filter (user_uuid);
CREATE INDEX idx_saved_filter_is_deleted ON saved_filter (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `shared_content_log`
-- --------------------------------------------------------
CREATE TABLE shared_content_log (
  shared_content_log_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  activity_uuid UUID DEFAULT NULL,
  message_queue_uuid UUID DEFAULT NULL,
  message_log_uuid UUID DEFAULT NULL,
  document_uuid UUID DEFAULT NULL,
  shared_url VARCHAR(1000) DEFAULT NULL,
  shared_title VARCHAR(255) DEFAULT NULL,
  recipient VARCHAR(255) DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (shared_content_log_uuid)
);

CREATE INDEX idx_shared_content_tenant ON shared_content_log (tenant_uuid);
CREATE INDEX idx_shared_content_workspace ON shared_content_log (workspace_uuid);
CREATE INDEX idx_shared_content_module_entity ON shared_content_log (module_uuid, entity_uuid);
CREATE INDEX idx_shared_content_activity ON shared_content_log (activity_uuid);
CREATE INDEX idx_shared_content_document ON shared_content_log (document_uuid);
CREATE INDEX idx_shared_content_queue ON shared_content_log (message_queue_uuid);
CREATE INDEX idx_shared_content_log_msg ON shared_content_log (message_log_uuid);
CREATE INDEX idx_shared_content_deleted ON shared_content_log (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `site_visit`
-- --------------------------------------------------------
CREATE TABLE site_visit (
  site_visit_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  opportunity_uuid UUID NOT NULL,
  activity_uuid UUID DEFAULT NULL,
  assigned_user_uuid UUID DEFAULT NULL,
  planned_visit_datetime TIMESTAMPTZ NOT NULL,
  actual_checkin_datetime TIMESTAMPTZ DEFAULT NULL,
  actual_checkout_datetime TIMESTAMPTZ DEFAULT NULL,
  reschedule_reason VARCHAR(1000) DEFAULT NULL,
  feedback_link VARCHAR(1000) DEFAULT NULL,
  rating DECIMAL(3,2) DEFAULT NULL,
  accompanying_members_details JSONB DEFAULT NULL,
  remarks VARCHAR(2000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (site_visit_uuid)
);

CREATE INDEX idx_site_visit_tenant ON site_visit (tenant_uuid);
CREATE INDEX idx_site_visit_workspace ON site_visit (workspace_uuid);
CREATE INDEX idx_site_visit_opportunity ON site_visit (opportunity_uuid);
CREATE INDEX idx_site_visit_activity ON site_visit (activity_uuid);
CREATE INDEX idx_site_visit_assigned_user ON site_visit (assigned_user_uuid);
CREATE INDEX idx_site_visit_planned_time ON site_visit (planned_visit_datetime);
CREATE INDEX idx_site_visit_checkin ON site_visit (actual_checkin_datetime);
CREATE INDEX idx_site_visit_deleted ON site_visit (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `storage`
-- --------------------------------------------------------
CREATE TABLE storage (
  storage_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  property_block_uuid UUID DEFAULT NULL,
  storage_number VARCHAR(100) NOT NULL,
  is_booked BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (storage_uuid),
  CONSTRAINT uk_workspace_storage_number UNIQUE (workspace_uuid, storage_number)
);

CREATE INDEX idx_storage_property_block_uuid ON storage (property_block_uuid);
CREATE INDEX idx_storage_is_deleted ON storage (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tag`
-- --------------------------------------------------------
CREATE TABLE tag (
  tag_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  tag_name VARCHAR(100) NOT NULL,
  is_system BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tag_uuid),
  CONSTRAINT uk_tag_name UNIQUE (tenant_uuid, workspace_uuid, module_uuid, tag_name)
);

CREATE INDEX idx_tag_tenant_uuid ON tag (tenant_uuid);
CREATE INDEX idx_tag_workspace_uuid ON tag (workspace_uuid);
CREATE INDEX idx_tag_is_archived ON tag (is_archived);
CREATE INDEX idx_tag_is_deleted ON tag (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tag_assigned`
-- --------------------------------------------------------
CREATE TABLE tag_assigned (
  tag_assigned_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  module_uuid UUID NOT NULL,
  entity_uuid UUID NOT NULL,
  tag_uuid UUID NOT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tag_assigned_uuid),
  CONSTRAINT uk_entity_tag UNIQUE (module_uuid, entity_uuid, tag_uuid)
);

CREATE INDEX idx_tag_assigned_tenant ON tag_assigned (tenant_uuid);
CREATE INDEX idx_tag_assigned_workspace ON tag_assigned (workspace_uuid);
CREATE INDEX idx_tag_assigned_module ON tag_assigned (module_uuid);
CREATE INDEX idx_tag_assigned_entity ON tag_assigned (entity_uuid);
CREATE INDEX idx_tag_assigned_tag ON tag_assigned (tag_uuid);
CREATE INDEX idx_tag_assigned_is_deleted ON tag_assigned (is_deleted);
CREATE INDEX idx_tag_assigned_module_entity ON tag_assigned (module_uuid, entity_uuid);

-- --------------------------------------------------------
-- Table structure for table `tenant`
-- --------------------------------------------------------
CREATE TABLE tenant (
  tenant_uuid UUID NOT NULL,
  tenant_code VARCHAR(50) NOT NULL,
  tenant_id VARCHAR(100) NOT NULL,
  tenant_name VARCHAR(255) NOT NULL,
  plan_type_universal_uuid UUID DEFAULT NULL,
  tenant_status_universal_uuid UUID DEFAULT NULL,
  contact_person_name VARCHAR(255) DEFAULT NULL,
  mobile_number VARCHAR(20) DEFAULT NULL,
  email_address VARCHAR(255) DEFAULT NULL,
  country_name VARCHAR(100) DEFAULT NULL,
  state_name VARCHAR(100) DEFAULT NULL,
  city_name VARCHAR(100) DEFAULT NULL,
  address TEXT DEFAULT NULL,
  gst_number VARCHAR(50) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_uuid),
  CONSTRAINT uk_tenant_code UNIQUE (tenant_code),
  CONSTRAINT uk_tenant_id UNIQUE (tenant_id)
);

CREATE INDEX idx_tenant_mobile_number ON tenant (mobile_number);
CREATE INDEX idx_tenant_email_address ON tenant (email_address);
CREATE INDEX idx_tenant_plan_type_universal_uuid ON tenant (plan_type_universal_uuid);
CREATE INDEX idx_tenant_status_universal_uuid ON tenant (tenant_status_universal_uuid);
CREATE INDEX idx_tenant_is_deleted ON tenant (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tenant_addon_purchase`
-- --------------------------------------------------------
CREATE TABLE tenant_addon_purchase (
  tenant_addon_purchase_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  addon_type_universal_uuid UUID DEFAULT NULL,
  quantity DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  purchase_date DATE NOT NULL,
  expiry_date DATE DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_addon_purchase_uuid)
);

CREATE INDEX idx_tenant_addon_purchase_tenant ON tenant_addon_purchase (tenant_uuid);
CREATE INDEX idx_tenant_addon_purchase_date ON tenant_addon_purchase (purchase_date);
CREATE INDEX idx_tenant_addon_purchase_type ON tenant_addon_purchase (addon_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `tenant_number_series`
-- --------------------------------------------------------
CREATE TABLE tenant_number_series (
  tenant_number_series_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID DEFAULT NULL,
  module_uuid UUID NOT NULL,
  prefix VARCHAR(50) NOT NULL,
  prefix_separator VARCHAR(10) DEFAULT NULL,
  current_number BIGINT NOT NULL DEFAULT 0,
  number_padding INT NOT NULL DEFAULT 6,
  reset_type INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_number_series_uuid),
  CONSTRAINT uk_tenant_workspace_entity UNIQUE (tenant_uuid, workspace_uuid, module_uuid)
);

CREATE INDEX idx_tenant_num_series_tenant ON tenant_number_series (tenant_uuid);
CREATE INDEX idx_tenant_num_series_workspace ON tenant_number_series (workspace_uuid);
CREATE INDEX idx_tenant_num_series_module ON tenant_number_series (module_uuid);
CREATE INDEX idx_tenant_num_series_deleted ON tenant_number_series (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tenant_setting`
-- --------------------------------------------------------
CREATE TABLE tenant_setting (
  tenant_setting_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  currency_uuid UUID NOT NULL,
  timezone_name VARCHAR(100) DEFAULT NULL,
  date_format VARCHAR(50) DEFAULT NULL,
  default_country_name VARCHAR(100) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_setting_uuid),
  CONSTRAINT uk_tenant_setting_tenant UNIQUE (tenant_uuid)
);

CREATE INDEX idx_tenant_setting_currency ON tenant_setting (currency_uuid);
CREATE INDEX idx_tenant_setting_deleted ON tenant_setting (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tenant_subscription`
-- --------------------------------------------------------
CREATE TABLE tenant_subscription (
  tenant_subscription_uuid UUID NOT NULL,
  subscription_number VARCHAR(50) NOT NULL,
  tenant_uuid UUID NOT NULL,
  plan_price_uuid UUID NOT NULL,
  renewal_of_subscription_uuid UUID DEFAULT NULL,
  subscription_start_date DATE NOT NULL,
  subscription_end_date DATE NOT NULL,
  grace_end_date DATE DEFAULT NULL,
  amount_paid DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_subscription_uuid),
  CONSTRAINT uk_subscription_number UNIQUE (subscription_number)
);

CREATE INDEX idx_tenant_sub_tenant ON tenant_subscription (tenant_uuid);
CREATE INDEX idx_tenant_sub_price ON tenant_subscription (plan_price_uuid);
CREATE INDEX idx_tenant_sub_end_date ON tenant_subscription (subscription_end_date);
CREATE INDEX idx_tenant_sub_deleted ON tenant_subscription (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tenant_subscription_payment`
-- --------------------------------------------------------
CREATE TABLE tenant_subscription_payment (
  tenant_subscription_payment_uuid UUID NOT NULL,
  payment_number VARCHAR(50) NOT NULL,
  tenant_uuid UUID NOT NULL,
  tenant_subscription_uuid UUID DEFAULT NULL,
  payment_mode_universal_uuid UUID DEFAULT NULL,
  payment_status_universal_uuid UUID DEFAULT NULL,
  payment_date DATE NOT NULL,
  amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  reference_number VARCHAR(100) DEFAULT NULL,
  receipt_file_name VARCHAR(255) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_subscription_payment_uuid),
  CONSTRAINT uk_payment_number UNIQUE (payment_number)
);

CREATE INDEX idx_tenant_sub_pay_tenant ON tenant_subscription_payment (tenant_uuid);
CREATE INDEX idx_tenant_sub_pay_sub ON tenant_subscription_payment (tenant_subscription_uuid);
CREATE INDEX idx_tenant_sub_pay_date ON tenant_subscription_payment (payment_date);
CREATE INDEX idx_tenant_sub_pay_mode ON tenant_subscription_payment (payment_mode_universal_uuid);
CREATE INDEX idx_tenant_sub_pay_status ON tenant_subscription_payment (payment_status_universal_uuid);
CREATE INDEX idx_tenant_sub_pay_deleted ON tenant_subscription_payment (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `tenant_usage_ledger`
-- --------------------------------------------------------
CREATE TABLE tenant_usage_ledger (
  tenant_usage_ledger_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  addon_type_universal_uuid UUID DEFAULT NULL,
  usage_transaction_type_universal_uuid UUID DEFAULT NULL,
  transaction_date TIMESTAMPTZ NOT NULL,
  quantity DECIMAL(18,2) NOT NULL DEFAULT 0.00,
  reference_number VARCHAR(100) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (tenant_usage_ledger_uuid)
);

CREATE INDEX idx_tenant_usage_ledger_tenant ON tenant_usage_ledger (tenant_uuid);
CREATE INDEX idx_tenant_usage_ledger_date ON tenant_usage_ledger (transaction_date);
CREATE INDEX idx_tenant_usage_ledger_ref ON tenant_usage_ledger (reference_number);
CREATE INDEX idx_tenant_usage_ledger_addon ON tenant_usage_ledger (addon_type_universal_uuid);
CREATE INDEX idx_tenant_usage_ledger_tx_type ON tenant_usage_ledger (usage_transaction_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `unit`
-- --------------------------------------------------------
CREATE TABLE unit (
  unit_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  property_block_uuid UUID NOT NULL,
  unit_status_universal_uuid UUID DEFAULT NULL,
  unit_type_universal_uuid UUID DEFAULT NULL,
  unit_number VARCHAR(50) NOT NULL,
  floor_number INT DEFAULT NULL,
  carpet_area DECIMAL(18,2) DEFAULT NULL,
  builtup_area DECIMAL(18,2) DEFAULT NULL,
  saleable_area DECIMAL(18,2) DEFAULT NULL,
  area_unit_universal_uuid UUID DEFAULT NULL,
  facing_details JSONB DEFAULT NULL,
  additional_area_details JSONB DEFAULT NULL,
  base_rate DECIMAL(18,2) DEFAULT NULL,
  additional_charge_details JSONB DEFAULT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (unit_uuid),
  CONSTRAINT uk_property_block_unit_number UNIQUE (property_block_uuid, unit_number)
);

CREATE INDEX idx_unit_tenant_uuid ON unit (tenant_uuid);
CREATE INDEX idx_unit_workspace_uuid ON unit (workspace_uuid);
CREATE INDEX idx_unit_property_block_uuid ON unit (property_block_uuid);
CREATE INDEX idx_unit_number ON unit (unit_number);
CREATE INDEX idx_unit_floor_number ON unit (floor_number);
CREATE INDEX idx_unit_is_deleted ON unit (is_deleted);
CREATE INDEX idx_unit_area_unit_universal_uuid ON unit (area_unit_universal_uuid);
CREATE INDEX idx_unit_status_universal_uuid ON unit (unit_status_universal_uuid);
CREATE INDEX idx_unit_type_universal_uuid ON unit (unit_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `unit_hold`
-- --------------------------------------------------------
CREATE TABLE unit_hold (
  unit_hold_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  unit_uuid UUID NOT NULL,
  contact_uuid UUID NOT NULL,
  held_by_user_uuid UUID NOT NULL,
  hold_from TIMESTAMPTZ NOT NULL,
  hold_until TIMESTAMPTZ NOT NULL,
  status_uuid UUID NOT NULL,
  remarks VARCHAR(1000) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  is_archived BOOLEAN NOT NULL DEFAULT FALSE,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (unit_hold_uuid)
);

CREATE INDEX idx_unit_hold_tenant_uuid ON unit_hold (tenant_uuid);
CREATE INDEX idx_unit_hold_workspace_uuid ON unit_hold (workspace_uuid);
CREATE INDEX idx_unit_hold_unit_uuid ON unit_hold (unit_uuid);
CREATE INDEX idx_unit_hold_contact_uuid ON unit_hold (contact_uuid);
CREATE INDEX idx_unit_hold_hold_until ON unit_hold (hold_until);

-- --------------------------------------------------------
-- Table structure for table `user`
-- --------------------------------------------------------
CREATE TABLE "user" (
  user_uuid UUID NOT NULL,
  user_status_universal_uuid UUID DEFAULT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) DEFAULT NULL,
  mobile_number VARCHAR(20) DEFAULT NULL,
  email_address VARCHAR(255) DEFAULT NULL,
  profile_photo_file_name VARCHAR(255) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_uuid)
);

CREATE INDEX idx_user_mobile_number ON "user" (mobile_number);
CREATE INDEX idx_user_email_address ON "user" (email_address);
CREATE INDEX idx_user_status_universal_uuid ON "user" (user_status_universal_uuid);
CREATE INDEX idx_user_is_deleted ON "user" (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `user_auth`
-- --------------------------------------------------------
CREATE TABLE user_auth (
  user_auth_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  login_id VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  is_mobile_verified BOOLEAN NOT NULL DEFAULT FALSE,
  is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  last_password_changed_at TIMESTAMPTZ DEFAULT NULL,
  last_login_at TIMESTAMPTZ DEFAULT NULL,
  failed_attempt_count INT NOT NULL DEFAULT 0,
  locked_until TIMESTAMPTZ DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_auth_uuid),
  CONSTRAINT uk_tenant_login_id UNIQUE (tenant_uuid, login_id)
);

CREATE INDEX idx_user_auth_user_uuid ON user_auth (user_uuid);
CREATE INDEX idx_user_auth_tenant_uuid ON user_auth (tenant_uuid);

-- --------------------------------------------------------
-- Table structure for table `user_device`
-- --------------------------------------------------------
CREATE TABLE user_device (
  user_device_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  device_type_universal_uuid UUID DEFAULT NULL,
  device_uuid VARCHAR(255) NOT NULL,
  device_name VARCHAR(255) DEFAULT NULL,
  app_version VARCHAR(50) DEFAULT NULL,
  device_os_version VARCHAR(100) DEFAULT NULL,
  firebase_token VARCHAR(500) DEFAULT NULL,
  last_seen_at TIMESTAMPTZ DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_device_uuid),
  CONSTRAINT uk_device_uuid UNIQUE (device_uuid)
);

CREATE INDEX idx_user_device_tenant_uuid ON user_device (tenant_uuid);
CREATE INDEX idx_user_device_user_uuid ON user_device (user_uuid);
CREATE INDEX idx_user_device_type_universal ON user_device (device_type_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `user_login_history`
-- --------------------------------------------------------
CREATE TABLE user_login_history (
  user_login_history_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  login_source_universal_uuid UUID DEFAULT NULL,
  device_uuid UUID DEFAULT NULL,
  login_datetime TIMESTAMPTZ NOT NULL,
  login_result_universal_uuid UUID DEFAULT NULL,
  ip_address VARCHAR(100) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_login_history_uuid)
);

CREATE INDEX idx_user_login_hist_tenant ON user_login_history (tenant_uuid);
CREATE INDEX idx_user_login_hist_user ON user_login_history (user_uuid);
CREATE INDEX idx_user_login_hist_datetime ON user_login_history (login_datetime);
CREATE INDEX idx_user_login_hist_result ON user_login_history (login_result_universal_uuid);
CREATE INDEX idx_user_login_hist_source ON user_login_history (login_source_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `user_otp`
-- --------------------------------------------------------
CREATE TABLE user_otp (
  user_otp_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_uuid UUID DEFAULT NULL,
  otp_purpose_universal_uuid UUID DEFAULT NULL,
  mobile_number VARCHAR(20) DEFAULT NULL,
  email_address VARCHAR(255) DEFAULT NULL,
  otp_code VARCHAR(20) NOT NULL,
  sent_at TIMESTAMPTZ DEFAULT NULL,
  attempt_count INT NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ NOT NULL,
  verified_at TIMESTAMPTZ DEFAULT NULL,
  is_consumed BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_otp_uuid)
);

CREATE INDEX idx_user_otp_tenant_uuid ON user_otp (tenant_uuid);
CREATE INDEX idx_user_otp_user_uuid ON user_otp (user_uuid);
CREATE INDEX idx_user_otp_expires_at ON user_otp (expires_at);
CREATE INDEX idx_user_otp_purpose_universal ON user_otp (otp_purpose_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `user_session`
-- --------------------------------------------------------
CREATE TABLE user_session (
  user_session_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  login_source_universal_uuid UUID DEFAULT NULL,
  device_uuid UUID DEFAULT NULL,
  jwt_id VARCHAR(255) DEFAULT NULL,
  login_at TIMESTAMPTZ NOT NULL,
  logout_at TIMESTAMPTZ DEFAULT NULL,
  expires_at TIMESTAMPTZ DEFAULT NULL,
  ip_address VARCHAR(100) DEFAULT NULL,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_session_uuid)
);

CREATE INDEX idx_user_session_tenant_uuid ON user_session (tenant_uuid);
CREATE INDEX idx_user_session_user_uuid ON user_session (user_uuid);
CREATE INDEX idx_user_session_device_uuid ON user_session (device_uuid);
CREATE INDEX idx_user_session_login_at ON user_session (login_at);
CREATE INDEX idx_user_session_source_universal ON user_session (login_source_universal_uuid);

-- --------------------------------------------------------
-- Table structure for table `user_tenant`
-- --------------------------------------------------------
CREATE TABLE user_tenant (
  user_tenant_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_uuid UUID NOT NULL,
  employee_code VARCHAR(50) DEFAULT NULL,
  designation VARCHAR(100) DEFAULT NULL,
  profile_details JSONB DEFAULT NULL,
  joining_date DATE DEFAULT NULL,
  leaving_date DATE DEFAULT NULL,
  is_default_tenant BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_tenant_uuid),
  CONSTRAINT uk_user_tenant UNIQUE (tenant_uuid, user_uuid)
);

CREATE INDEX idx_user_tenant_user_uuid ON user_tenant (user_uuid);
CREATE INDEX idx_user_tenant_tenant_uuid ON user_tenant (tenant_uuid);
CREATE INDEX idx_user_tenant_is_deleted ON user_tenant (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `user_workspace`
-- --------------------------------------------------------
CREATE TABLE user_workspace (
  user_workspace_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  user_tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  role_uuid UUID NOT NULL,
  assigned_at TIMESTAMPTZ DEFAULT NULL,
  is_default_workspace BOOLEAN NOT NULL DEFAULT FALSE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (user_workspace_uuid),
  CONSTRAINT uk_user_workspace UNIQUE (user_tenant_uuid, workspace_uuid)
);

CREATE INDEX idx_user_workspace_tenant_uuid ON user_workspace (tenant_uuid);
CREATE INDEX idx_user_workspace_workspace_uuid ON user_workspace (workspace_uuid);
CREATE INDEX idx_user_workspace_role_uuid ON user_workspace (role_uuid);
CREATE INDEX idx_user_workspace_is_deleted ON user_workspace (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `webhook_log`
-- --------------------------------------------------------
CREATE TABLE webhook_log (
  webhook_log_uuid UUID NOT NULL,
  tenant_uuid UUID DEFAULT NULL,
  source_name VARCHAR(100) NOT NULL,
  event_name VARCHAR(150) NOT NULL,
  external_reference VARCHAR(255) DEFAULT NULL,
  request_payload_json JSONB DEFAULT NULL,
  response_payload_json JSONB DEFAULT NULL,
  processed_at TIMESTAMPTZ DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  PRIMARY KEY (webhook_log_uuid)
);

-- --------------------------------------------------------
-- Table structure for table `workspace`
-- --------------------------------------------------------
CREATE TABLE workspace (
  workspace_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_name VARCHAR(255) NOT NULL,
  workspace_code VARCHAR(100) DEFAULT NULL,
  workspace_prefix VARCHAR(20) DEFAULT NULL,
  public_slug VARCHAR(150) DEFAULT NULL,
  microsite_settings JSONB DEFAULT NULL,
  default_owner_user_uuid UUID DEFAULT NULL,
  is_open_workspace BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (workspace_uuid),
  CONSTRAINT uk_tenant_workspace_name UNIQUE (tenant_uuid, workspace_name),
  CONSTRAINT uk_tenant_public_slug UNIQUE (tenant_uuid, public_slug)
);

CREATE INDEX idx_workspace_tenant_uuid ON workspace (tenant_uuid);
CREATE INDEX idx_workspace_default_owner_user_uuid ON workspace (default_owner_user_uuid);
CREATE INDEX idx_workspace_is_deleted ON workspace (is_deleted);

-- --------------------------------------------------------
-- Table structure for table `workspace_assignment_rule`
-- --------------------------------------------------------
CREATE TABLE workspace_assignment_rule (
  workspace_assignment_rule_uuid UUID NOT NULL,
  tenant_uuid UUID NOT NULL,
  workspace_uuid UUID NOT NULL,
  assignment_type_universal_uuid UUID DEFAULT NULL,
  default_role_uuid UUID DEFAULT NULL,
  default_user_uuid UUID DEFAULT NULL,
  last_assigned_user_uuid UUID DEFAULT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  remarks VARCHAR(500) DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NULL,
  created_by_user_uuid UUID DEFAULT NULL,
  updated_at TIMESTAMPTZ DEFAULT NULL,
  updated_by_user_uuid UUID DEFAULT NULL,
  archived_at TIMESTAMPTZ DEFAULT NULL,
  archived_by_user_uuid UUID DEFAULT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  deleted_by_user_uuid UUID DEFAULT NULL,
  PRIMARY KEY (workspace_assignment_rule_uuid),
  CONSTRAINT uk_workspace_assignment UNIQUE (tenant_uuid, workspace_uuid)
);

CREATE INDEX idx_workspace_assignment_rule_tenant ON workspace_assignment_rule (tenant_uuid);
CREATE INDEX idx_workspace_assignment_rule_workspace ON workspace_assignment_rule (workspace_uuid);
CREATE INDEX idx_workspace_assignment_rule_assignment_type ON workspace_assignment_rule (assignment_type_universal_uuid);
CREATE INDEX idx_workspace_assignment_rule_is_deleted ON workspace_assignment_rule (is_deleted);