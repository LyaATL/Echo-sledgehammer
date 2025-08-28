# Echo sledgehammer

A simple centralized server to sync and update a global ban list.
People that host their own echo-service instance will be able to submit a ban made on their own service to be included in the global list.

## Purpose

The goal of this project is to maintain a trusted, shared list of banned users across multiple echo-service instances.
This reduces the risk of harmful actors simply moving from one service to another after being banned locally.
Sledgehammer wil also inform echo-service instances of new/fresh players, to easily inform administrators/moderators.

## Criteria for Global Ban

A global ban is only applied in cases of severe and verified misconduct.
While local bans may vary depending on a host’s rules, inclusion in the global ban list is reserved for actions such as (but not limited to):

- Child exploitation or pedophilia (e.g., possession, sharing, or promotion of such material)
- Sexual assault, harassment, or grooming
- Terrorism-related activity (organizing, recruiting, or promoting violent extremism)
- Severe threats of violence or credible intent to harm others
- Distribution of malware, botnets, or large-scale attacks against other services
- Identity theft, fraud, or large-scale scams
- Organized hate speech (targeted campaigns based on race, religion, gender, sexuality, etc.)

These bans are intended to protect the entire community from individuals who pose a significant risk across platforms.

## How It Works

Each host can ban users locally according to their own rules.
When a ban meets the global criteria, the host may submit it to the centralized server.

Submissions are verified and, if approved, propagated to the global ban list.
Other echo-service instances can sync this list to keep their communities safe.

## Notes
Transparency: Ban reasons should be clearly documented when submitted.
Appeals: If appropriate, processes for disputing or reviewing global bans should be implemented.
Security: The system should prevent abuse of the submission mechanism (e.g., requiring signatures, moderator authentication, or quorum approval).
