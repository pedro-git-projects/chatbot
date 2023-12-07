ALTER TABLE users ADD CONSTRAINT valid_role CHECK (role IN ('admin', 'collaborator', 'user')); 
