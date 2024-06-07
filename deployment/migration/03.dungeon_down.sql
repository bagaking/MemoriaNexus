-- Dropping the dungeon-related tables. The order of dropping is reverse of creation.
DROP TABLE IF EXISTS `dungeon_monsters`;
DROP TABLE IF EXISTS `dungeon_items`;
DROP TABLE IF EXISTS `dungeon_tags`;
DROP TABLE IF EXISTS `dungeon_books`;
DROP TABLE IF EXISTS `user_monsters`;
DROP TABLE IF EXISTS `monsters`;
DROP TABLE IF EXISTS `dungeons`;
-- Tables are now dropped, and the database schema is returned to the state before the up migration script was applied.