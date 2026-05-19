-- 1. Trigger on contacts (For Insert/Update)
-- This prevents saving a contact with an invalid is_payment_responsible state.
CREATE OR REPLACE FUNCTION validate_contact_payment_responsible()
RETURNS TRIGGER AS $$
DECLARE
    v_is_organic BOOLEAN;
BEGIN
    -- Only check if the contact is marked as payment responsible
    IF NEW.is_payment_responsible = TRUE THEN
        
        -- Scenario A: Contact belongs directly to an organization
        IF NEW.organization_id IS NOT NULL THEN
            SELECT is_organic INTO v_is_organic FROM organizations WHERE id = NEW.organization_id;
            IF NOT v_is_organic THEN
                RAISE EXCEPTION 'A contact belonging to an organization can only be payment responsible if the organization is organic.';
            END IF;
            
        -- Scenario B: Contact belongs to an office
        ELSIF NEW.office_id IS NOT NULL THEN
            SELECT org.is_organic INTO v_is_organic
            FROM offices off
            JOIN organizations org ON org.id = off.organization_id
            WHERE off.id = NEW.office_id;
            
            IF v_is_organic THEN
                RAISE EXCEPTION 'A contact belonging to an office can only be payment responsible if the organization is inorganic.';
            END IF;
        END IF;

    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_contact_payment_responsible
BEFORE INSERT OR UPDATE OF is_payment_responsible, office_id, organization_id ON contacts
FOR EACH ROW
EXECUTE PROCEDURE validate_contact_payment_responsible();

-- 2. Trigger on organizations (When is_organic changes)
-- This blocks an organization update if it would leave the database in an invalid state.
CREATE OR REPLACE FUNCTION validate_org_organic_change()
RETURNS TRIGGER AS $$
BEGIN
    -- We only care if the is_organic field actually changes
    IF NEW.is_organic <> OLD.is_organic THEN
        
        -- If changing to Organic: No offices belonging to this org can have payment responsible contacts
        IF NEW.is_organic = TRUE THEN
            IF EXISTS (
                SELECT 1 FROM contacts c
                JOIN offices o ON o.id = c.office_id
                WHERE o.organization_id = NEW.id AND c.is_payment_responsible = TRUE
            ) THEN
                RAISE EXCEPTION 'Cannot change organization to Organic: one or more of its offices have a payment responsible contact.';
            END IF;
            
        -- If changing to Inorganic: The organization itself cannot have payment responsible contacts
        ELSE
            IF EXISTS (
                SELECT 1 FROM contacts c
                WHERE c.organization_id = NEW.id AND c.is_payment_responsible = TRUE
            ) THEN
                RAISE EXCEPTION 'Cannot change organization to Inorganic: the organization itself has a payment responsible contact directly attached.';
            END IF;
        END IF;
        
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_org_organic_change
BEFORE UPDATE OF is_organic ON organizations
FOR EACH ROW
EXECUTE PROCEDURE validate_org_organic_change();

-- 3. Trigger on offices (When organization_id changes)
-- Prevents moving an office with a payment responsible contact to an organic organization.
CREATE OR REPLACE FUNCTION validate_office_org_change()
RETURNS TRIGGER AS $$
DECLARE
    v_new_org_is_organic BOOLEAN;
BEGIN
    IF NEW.organization_id <> OLD.organization_id THEN
        SELECT is_organic INTO v_new_org_is_organic FROM organizations WHERE id = NEW.organization_id;
        
        IF v_new_org_is_organic = TRUE AND EXISTS (
            SELECT 1 FROM contacts WHERE office_id = NEW.id AND is_payment_responsible = TRUE
        ) THEN
             RAISE EXCEPTION 'Cannot move office to an Organic organization: the office already has a payment responsible contact.';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_office_org_change
BEFORE UPDATE OF organization_id ON offices
FOR EACH ROW
EXECUTE PROCEDURE validate_office_org_change();
