namespace Backend.Services;

using System;
using System.Collections.Generic;
using Microsoft.Data.Sqlite;

public class UserDatabase
{
    private readonly string _connectionString;

    public UserDatabase(string dbPath)
    {
        _connectionString = $"Data Source={dbPath}";
        InitializeDatabase();
    }

    private void InitializeDatabase()
    {
        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = @"
                CREATE TABLE IF NOT EXISTS Users (
                    Uid TEXT PRIMARY KEY,
                    Name TEXT NOT NULL,
                    Pin TEXT NOT NULL
                );

                CREATE TABLE IF NOT EXISTS LoginRecords (
                    Id INTEGER PRIMARY KEY AUTOINCREMENT,
                    Uid TEXT NOT NULL,
                    IP TEXT NOT NULL,
                    Browser TEXT NOT NULL,
                    LoginTime TEXT NOT NULL,
                    FOREIGN KEY (Uid) REFERENCES Users(Uid)
                );";
            command.ExecuteNonQuery();
        }
    }

    public void AddUser(string uid, string name, string pin)
    {
        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = "INSERT INTO Users (Uid, Name, Pin) VALUES (@uid, @name, @pin)";
            command.Parameters.AddWithValue("@uid", uid);
            command.Parameters.AddWithValue("@name", name);
            command.Parameters.AddWithValue("@pin", pin);

            command.ExecuteNonQuery();
        }
    }

    public List<string> ValidateUser(string name, string pin)
    {
        var validUids = new List<string>();

        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = "SELECT Uid FROM Users WHERE Name = @name AND Pin = @pin";
            command.Parameters.AddWithValue("@name", name);
            command.Parameters.AddWithValue("@pin", pin);

            using (var reader = command.ExecuteReader())
            {
                while (reader.Read())
                {
                    validUids.Add(reader.GetString(0));
                }
            }
        }

        return validUids;
    }

    public void AddLoginRecord(string uid, string ip, string browser)
    {
        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = @"
                INSERT INTO LoginRecords (Uid, IP, Browser, LoginTime) 
                VALUES (@uid, @ip, @browser, @loginTime)";
            command.Parameters.AddWithValue("@uid", uid);
            command.Parameters.AddWithValue("@ip", ip);
            command.Parameters.AddWithValue("@browser", browser);
            command.Parameters.AddWithValue("@loginTime", DateTime.UtcNow.ToString("o"));

            command.ExecuteNonQuery();
        }
    }

    public List<(string IP, string Browser, DateTime LoginTime)> GetUserLoginHistory(string uid)
    {
        var history = new List<(string IP, string Browser, DateTime LoginTime)>();

        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = "SELECT IP, Browser, LoginTime FROM LoginRecords WHERE Uid = @uid ORDER BY LoginTime DESC";
            command.Parameters.AddWithValue("@uid", uid);

            using (var reader = command.ExecuteReader())
            {
                while (reader.Read())
                {
                    history.Add((
                        reader.GetString(0),
                        reader.GetString(1),
                        DateTime.Parse(reader.GetString(2))
                    ));
                }
            }
        }

        return history;
    }

    public bool UserExists(string uid)
    {
        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = "SELECT COUNT(*) FROM Users WHERE Uid = @uid";
            command.Parameters.AddWithValue("@uid", uid);

            var count = Convert.ToInt32(command.ExecuteScalar());
            return count > 0;
        }
    }

    public bool UpdateUser(string uid, string newName, string newPin)
    {
        using (var connection = new SqliteConnection(_connectionString))
        {
            connection.Open();

            var command = connection.CreateCommand();
            command.CommandText = "UPDATE Users SET Name = @name, Pin = @pin WHERE Uid = @uid";
            command.Parameters.AddWithValue("@uid", uid);
            command.Parameters.AddWithValue("@name", newName);
            command.Parameters.AddWithValue("@pin", newPin);

            int rowsAffected = command.ExecuteNonQuery();
            return rowsAffected > 0;
        }
    }
}