REST API Endpoints
- /users (GET, POST, PUT)
    - ID
    - Profile image url
    - Name
    - Time zone
    - Phone number
    - Created At Date & Time
- /habits (GET, POST, PUT)
    - User ID (references a User)
    - ID
    - Name
    - Description
    - Created At Date & Time
- /logs (GET, POST, PUT)
    - Habit ID (references a Habit)
    - Created At Date & Time
    - Notes

Create a UI using react to interact with the above API that will be hosting it
this UI helps users keep track of habits in streaks
Each time a habit log is created it adds to the streak
The users goal is to keep their habit streak going as long as possible without breaking it

The style should be modern and minimal utilizing sans serif fonts, pastel color schemes, and rounded corners

UI Component structure
    - HabitList
        - Habit
            - Name
            - Description
            - DeleteButton
                - deletes habit using API
            - EditButton
                - opens HabitDetailsModal
            - StreakCount
                - The user is allowed to skip a day before their streak count resets
            - NewLogButton (creates log for current date)
                - opens LogDetailsModal to fill in details
            - StreakList
                - DayContainer
                    - if this day has a habit log
                        - Habit Log (Edit on click opening the same modal as before)
                    - else
                        - empty
            - LogDetailsModal (hidden by default)
                - NotesField
                - SaveButton
                    - Inserts / updates log using API
    - AddNewHabitButton
        - opens HabitDetailsModal
    - HabitDetailsModal (hidden by default)
        - NameField
        - DescriptionField
        - SaveButton
            - inserts / updates habit using API


- all errors are gracefully handled and displayed
- Use standard libraries where possible to avoid excessive dependencies
- modular components and functions
- Robust, idiomatic, clean code 
- UI is well formatted and all components are properly sized on screen with no run off
- clear comments on all functions and conditionals
- All buttons perform the correct action passing data to API
- All data is properly fetched from the api
- README.md file with instructions to run locally and to deploy to AWS 