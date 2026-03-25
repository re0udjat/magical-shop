ALTER TABLE items ADD CONSTRAINT name_check CHECK (LENGTH(name) BETWEEN 1 AND 500);
ALTER TABLE items ADD CONSTRAINT rarity_check CHECK (rarity IN ('common', 'uncommon', 'rare', 'mythic', 'legendary'));
ALTER TABLE items ADD CONSTRAINT price_check CHECK (price > 0 AND price <= 1_000_000_000);