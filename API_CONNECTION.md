# Connecting to the Backend API

The React Native app needs to connect to the Go backend server. By default, it's configured to use `http://localhost:8080`.

## Different Platform Configurations

### iOS Simulator
- Uses `localhost` or `127.0.0.1`
- Default configuration works: `http://localhost:8080`

### Android Emulator
- Cannot use `localhost` (it refers to the emulator itself)
- Use `10.0.2.2` to access host machine
- Update `API_URL` in `src/services/api.ts`:
  ```typescript
  const API_URL = 'http://10.0.2.2:8080';
  ```

### Physical Device
- Both iOS and Android need your computer's IP address
- Find your IP address:
  - Mac/Linux: `ifconfig | grep "inet "`
  - Windows: `ipconfig`
- Update `API_URL` in `src/services/api.ts`:
  ```typescript
  const API_URL = 'http://YOUR_IP_ADDRESS:8080';
  ```
- Make sure your phone and computer are on the same WiFi network

## Making the Change

1. Open `src/services/api.ts`
2. Find the line: `const API_URL = 'http://localhost:8080';`
3. Replace with the appropriate URL for your platform
4. Restart the Metro bundler and reload the app
