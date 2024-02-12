-- Down migration scripts for [[Profile]] Module

-- Removing the tables created by the 01.passport_up.sql migration
-- migration/01.passport_down.sql

-- Since we didn't create foreign key constraints, table dropping order doesn't matter for constraints,
-- but it's still good practice to reverse the creation order.

-- Dropping the two_factor_settings table
DROP TABLE IF EXISTS `user_two_factor_settings`;

-- Dropping the user_oauth_credentials table
DROP TABLE IF EXISTS `user_oauth_credentials`;

-- Dropping the users table
DROP TABLE IF EXISTS `users`;

-- All the associated tables for the Profile Module's passport feature are now successfully removed.