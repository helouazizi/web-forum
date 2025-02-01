#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Open SQLite database
sqlite3* open_database(const char *filename) {
    sqlite3 *db;
    if (sqlite3_open(filename, &db) != SQLITE_OK) {
        return NULL;  // Failed to open DB
    }
    return db;
}

// Execute an SQL statement (e.g., create table, insert data)
int execute_sql(sqlite3 *db, const char *sql) {
    char *err_msg = 0;
    int rc = sqlite3_exec(db, sql, 0, 0, &err_msg);
    if (rc != SQLITE_OK) {
        sqlite3_free(err_msg);
        return rc;
    }
    return SQLITE_OK;
}

// Retrieve data from a table
void fetch_users(sqlite3 *db) {
    sqlite3_stmt *stmt;
    const char *sql = "SELECT id, name FROM users;";

    if (sqlite3_prepare_v2(db, sql, -1, &stmt, 0) != SQLITE_OK) {
        return;  // Failed to prepare SQL
    }

    while (sqlite3_step(stmt) == SQLITE_ROW) {
        int id = sqlite3_column_int(stmt, 0);
        const char *name = (const char*)sqlite3_column_text(stmt, 1);
        printf("User ID: %d, Name: %s\n", id, name);
    }

    sqlite3_finalize(stmt);
}

// Close the database
void close_database(sqlite3 *db) {
    sqlite3_close(db);
}
