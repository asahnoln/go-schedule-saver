# About

This app is a parser which takes a source with events and saves to a destination.

## Example of Source

| Time / Days | Jan 1, Mon | Jan 3, Wed |                                               |
| ----------- | ---------- | ---------- | --------------------------------------------- |
| 12:00       |            | Ivan       | Extra info for the table (shouldn't be saved) |
| 12:30       | John       |            |                                               |
| 13:00       | Michael    |            |                                               |

# TODO

- [x] Source, which parses information and returns data.
- [x] Destination, which takes data from Source and saves somewhere.
- [x] Source: link to a Google Spreadsheet
- [x] Cmd with destination in memory
- [x] Destination: Google Calendar

## Happy Path

- [x] Function, which takes Source, Destination, and returns err|nil. It should take info from Source and save it to Destination

## Unhappy Path

- [x] Source returns error
- [x] Destination returns error
- [x] Saver returns error
- [ ] Source, Destination return different type errors
- [x] Google Sheets - error handling when creating
- [x] Google Sheets - handle edge cases
- [ ] Google Calendar check for existing events and don't save them
- [ ] Google Calendar remove existing events if their name is duplicated at different time
