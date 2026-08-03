-- Enable pgcrypto for uuid generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$ 
DECLARE
    v_tenant_uuid UUID;
    v_sub_uuid UUID;
    v_user_uuid UUID;
    
    v_status_active UUID;
    v_payment_paid UUID;
    v_plan_price UUID;
    
    i INT;
    j INT;
BEGIN
    -- Fetch active statuses and a plan price to reuse
    SELECT universal_uuid INTO v_status_active FROM mst_universal WHERE universal_name ILIKE '%Active%' LIMIT 1;
    SELECT universal_uuid INTO v_payment_paid FROM mst_universal WHERE universal_name ILIKE '%Paid%' LIMIT 1;
    SELECT plan_price_uuid INTO v_plan_price FROM mst_plan_price LIMIT 1;

    -- Generate 50 Tenants
    FOR i IN 1..50 LOOP
        v_tenant_uuid := gen_random_uuid();
        
        INSERT INTO tenant (tenant_uuid, tenant_code, tenant_id, tenant_name, tenant_status_universal_uuid, contact_person_name, mobile_number, email_address, created_at, updated_at)
        VALUES (
            v_tenant_uuid, 
            'T-' || lpad(i::text, 4, '0'), 
            'T-' || lpad(i::text, 4, '0'), 
            'Mock Tenant Corp ' || i, 
            v_status_active, 
            'Admin ' || i, 
            '98000' || lpad(i::text, 5, '0'), 
            'admin' || i || '@mocktenant.com', 
            NOW() - (random() * 365 * interval '1 day'), 
            NOW()
        );

        -- Add 1 Subscription for the tenant
        IF v_plan_price IS NOT NULL THEN
            v_sub_uuid := gen_random_uuid();
            INSERT INTO tenant_subscription (tenant_subscription_uuid, subscription_number, tenant_uuid, plan_price_uuid, subscription_start_date, subscription_end_date, amount_paid, created_at, updated_at)
            VALUES (
                v_sub_uuid, 
                'SUB-2026-' || lpad(i::text, 4, '0'), 
                v_tenant_uuid, 
                v_plan_price, 
                CURRENT_DATE - (random() * 100 * interval '1 day')::interval, 
                CURRENT_DATE + (random() * 200 * interval '1 day')::interval, 
                9999.00, 
                NOW(), 
                NOW()
            );
            
            -- Add 1 Payment for the subscription
            INSERT INTO tenant_subscription_payment (tenant_subscription_payment_uuid, payment_number, tenant_uuid, tenant_subscription_uuid, payment_status_universal_uuid, payment_date, amount, created_at)
            VALUES (
                gen_random_uuid(), 
                'PAY-' || lpad(i::text, 5, '0'), 
                v_tenant_uuid, 
                v_sub_uuid, 
                v_payment_paid, 
                CURRENT_DATE - (random() * 30 * interval '1 day')::interval, 
                9999.00, 
                NOW()
            );
        END IF;

        -- Generate 50 Users for each Tenant
        FOR j IN 1..50 LOOP
            v_user_uuid := gen_random_uuid();
            
            -- User Profile
            INSERT INTO "user" (user_uuid, user_status_universal_uuid, first_name, last_name, email_address, created_at, updated_at)
            VALUES (
                v_user_uuid, 
                v_status_active, 
                'User' || j, 
                'Tenant' || i, 
                'user' || j || '@mocktenant' || i || '.com', 
                NOW(), 
                NOW()
            );

            -- User Auth mapping to Tenant
            INSERT INTO user_auth (user_auth_uuid, tenant_uuid, user_uuid, login_id, password_hash, created_at, updated_at)
            VALUES (
                gen_random_uuid(), 
                v_tenant_uuid, 
                v_user_uuid, 
                'user' || j || '@mocktenant' || i || '.com', 
                'mock_password_hash', 
                NOW(), 
                NOW()
            );
            
        END LOOP;
        
    END LOOP;

END $$;
