# Echo sledgehammer

A simple centralized server to sync and update a global ban list.
People that host their own echo-service instance will be able to submit a ban made on their own service to be included in the global list.

## Purpose

The goal of this project is to maintain a trusted, shared list of banned users across multiple echo-service instances.
This reduces the risk of harmful actors simply moving from one service to another after being banned locally.
Sledgehammer wil also inform echo-service instances of new/fresh players, to easily inform administrators/moderators.

## Criteria for Global client ban

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

## Global File Blacklist

In addition to the global ban list for users, Echo Sledgehammer also maintains a **global file blacklist**.  
This list contains file hashes and signatures of texture files (or similar modifications) that are known to cause crashes, instability, or other disruptive behavior when loaded in-game.  

Examples include:
- Texture or model files that are intentionally crafted to crash the client or overload resources  
- Files that cause rendering glitches or instability across multiple systems  
- Files distributed under the guise of normal mods but known to break core game functions  
- Mods or files that depict or promote **underage NSFW content**  
- Rare cases where a file may trigger a **0-day vulnerability** in the client. While extremely unlikely, the blacklist system provides a safeguard if such an issue is ever discovered 

Each file entry includes:  
- A cryptographic hash (for reliable identification)  
- The classification or reason for blacklisting (e.g., crash-texture, instability, exploit-suspected)  
- (Optional) reference notes for transparency  

Hosts can sync this blacklist and automatically block or warn against the use, upload, or distribution of flagged files within their own echo-service environment.  

This helps prevent the spread of disruptive or malicious texture files across the ecosystem and improves overall stability for players.  

## How It Works

Each host can ban users locally according to their own rules.
When a ban meets the global criteria, the host may submit it to the centralized server.

Submissions are verified and, if approved, propagated to the global ban list.
Other echo-service instances can sync this list to keep their communities safe.

## Privacy & Data Collection Disclaimer

The Echo sledgehammer system requires a way to map and match players to enforce bans consistently across services. This mapping is done by the client plugin (running within FFXIV via Dalamud) sending identifying information to an echo-service instance.

### What Information Is Collected

To remain as privacy-neutral as possible, we avoid gathering sensitive or personally identifying data. Instead, we rely on in-game identifiers that are already public in FFXIV. Examples include:
- Character Name (as displayed in-game)
- World / Data Center (to disambiguate players with the same name)
- Character ID / Lodestone ID (publicly available numeric ID tied to Lodestone profiles)

### What Information Is Not Collected

We do not collect or transmit:
- Hardware identifiers (HWID, MAC, serial numbers, etc.)
- Email addresses or account credentials
- Chat logs or personal messages
- Any other personal data outside the game context

## Privacy Commitments

Data sent to/from the Echo service, Echo mesh or Echo Sledgehamer is the minimum necessary to uniquely identify a character for ban purposes.<br/>
All ban entries must include the reason for the ban, but no additional personal data.<br/>
This ensures the global ban system remains effective without compromising player privacy.<br/>

## How to Run

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd echo-sledgehammer
   ```

2. **Start the service**
   ```bash
   docker-compose up --build
   ```

3. **Access the service**
    - API: http://localhost:8080
    - Metrics: http://localhost:8080/metrics

4. **Stop the service**
   ```bash
   docker-compose down
   ```

### Using Docker

1. **Build the image**
   ```bash
   docker build -t sledgehammer .
   ```

2. **Run the container**
   ```bash
   docker run -d \
     --name sledgehammer \
     -p 8080:8080 \
     -v sledgehammer_data:/data \
     -e DATABASE_PATH=/data/sledgehammer.db \
     sledgehammer
   ```

### Local Development

1. **Prerequisites**
    - Go 1.24 or later
    - SQLite3 development libraries

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export DATABASE_PATH=./sledgehammer.db
   ```

4. **Run the application**
   ```bash
   go run ./cmd/sledgehammer
   ```

### API Endpoints

- `GET /bans` - List all bans
- `POST /bans` - Submit a new ban
- `GET /metrics` - Prometheus metrics

### Environment Variables

- `DATABASE_PATH` - Path to SQLite database file (default: uses `os.Getenv("DATABASE_PATH")`)

## Notes
Transparency: Ban reasons should be clearly documented when submitted.<br/>
Appeals: If appropriate, processes for disputing or reviewing global bans should be implemented.<br/>

## Disclaimer

This project, **Echo Sledgehammer**, is an independent work and has **no affiliation, connection, or association with Mare Synchronos** or its developers.  
It does not use, depend on, or share any code or intellectual property from Mare Synchronos.

The Echo Sledgehammer ban list and file blacklist are provided **as-is**, without any guarantees or warranties.  
While every effort is made to ensure entries are accurate and justified, there is no guarantee that the lists are complete, error-free, or up-to-date.  

By using these lists, each host accepts full responsibility for how the information is applied within their own service.  
Echo Sledgehammer maintainers are not liable for any damages, disruptions, or disputes that may arise from their use.  

