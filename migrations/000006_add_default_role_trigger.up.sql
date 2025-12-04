CREATE OR REPLACE FUNCTION assign_default_role()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_roles (user_id, role_id)
    SELECT NEW.id, id FROM roles WHERE name = 'USER';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_default_role_trigger
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION assign_default_role();
