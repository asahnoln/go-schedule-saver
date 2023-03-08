# About

This app is a parser which takes a source with events and saves to a destination.

## Example of Source

| Time / Days | Jan 1, Mon | Jan 3, Wed |
| ----------- | ---------- | ---------- |
| 12:00       | x          | Ivan       |
| 12:30       | John       | x          |
| 13:00       | Michael    | x          |

# TODO

- [ ] Source, which parses information and returns data.
- [ ] Destination, which takes data from Source and saves somewhere.
- [ ] Source: link to a Google Spreadsheet
- [ ] Destination: Google Calendar

## Happy Path

- [x] Function, which takes Source, Destination, and returns err|nil. It should take info from Source and save it to Destination

## Unhappy Path

- [x] Source returns error
- [x] Destination returns error
- [x] Saver returns error
- [ ] Source, Destination return different type errors

