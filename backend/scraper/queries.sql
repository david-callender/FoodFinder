-- THROUGHOUT THIS FILE, JAN 8 1999 WILL BE USED AS A PLACEHOLDER DATE

-- Obtain all food items for a given date.
SELECT * FROM DocCache WHERE day='1999-01-08';

-- Prune all food items for menu dates before a given date.
DELETE FROM DocCache WHERE day<'1999-01-08';

-- Prune all food items for menu dates equal to a given date.
-- May be used to update item cache 
DELETE FROM DocCache WHERE day='1999-01-08' AND location='<LOCATION ID>' AND mealtime='<MEALTIME>';

-- Insert a new set of food items.
INSERT INTO DocCache (day, location, mealtime, meal, mealid)
	VALUES
	('1999-01-08', '<ID>', '<MEALTIME NAME>', '<MEAL NAME>', '<MEAL ID>');
	-- And so on...

