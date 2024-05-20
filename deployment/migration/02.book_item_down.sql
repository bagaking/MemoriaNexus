-- Down migration scripts to undo the profile-related schema changes
-- Dropping all the tables that were created in the up migration script. The order of dropping is reverse of creation.
-- Since there are no foreign keys, we do not need to be concerned about the order regarding constraints, but it is a good practice to reverse the creation order.

DROP TABLE IF EXISTS `book_items`;
DROP TABLE IF EXISTS `item_tags`;
DROP TABLE IF EXISTS `book_tags`;
DROP TABLE IF EXISTS `books`;
DROP TABLE IF EXISTS `items`;
DROP TABLE IF EXISTS `tags`;
-- Tables are now dropped, and the database schema is returned to the state before the up migration script was applied.