ALTER TABLE items ADD CONSTRAINTS name_check CHECK (LENGTH(name) BETWEEN 1 AND 500);
ALTER TABLE items ADD CONSTRAINTS rarity_check CHECK (rarity IN ('common', 'uncommon', 'rare', 'mythic', 'legendary'));
ALTER TABLE items ADD CONSTRAINTS price_check CHECK (price > 0 AND price <= 1_000_000_000);