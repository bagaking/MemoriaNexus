-- Down migration scripts to undo the profile-related schema changes
-- Dropping all the tables that were created in the up migration script. The order of dropping is reverse of creation.
-- Since there are no foreign keys, we do not need to be concerned about the order regarding constraints, but it is a good practice to reverse the creation order.

DROP TABLE IF EXISTS `profile`;
DROP TABLE IF EXISTS `profile_settings_memorization`;
DROP TABLE IF EXISTS `profile_settings_advance`;
DROP TABLE IF EXISTS `profile_points`;
-- Tables are now dropped, and the database schema is returned to the state before the up migration script was applied.