-- Select all matches between user preferences and the DocCache table
SELECT user, meal 
	FROM "Preferences"
	JOIN "DocCache" 
	ON "Preferences.preference" = "DocCache.meal" 
	WHERE day="2025-03-03"
	ORDER BY user;

-- Select emails from the users table where the user is in the preferences table
-- We would then build a map from uid to email from this table
SELECT email, id FROM "Users" JOIN "Preferences" ON "Users.id" = "Preferences.user";
