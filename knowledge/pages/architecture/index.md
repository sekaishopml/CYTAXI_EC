# Architecture

## System overview
- [[frontend/index]]
- [[backend/index]]
- [[packages/index]]

## Data flow
1. User opens miniweb → state machine initializes at `pickup_select`
2. Map detects GPS → reverse geocode → pickup selection
3. User confirms pickup → search destination → select → confirm
4. Trip request → driver matching → tracking → payment → rating
