# NGINX Config Generator

A web-based tool for generating NGINX configuration files easily.

## Quick Start

1. Clone the repository and navigate to the project directory:
   ```
   git clone https://github.com/PatrickSUDO/nginx-config.git
   cd nginx-config
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Start the server:
   ```
   go run main.go
   ```

4. Open your web browser and go to `http://localhost:8080`

## How to Use

### Option 1: Upload a YAML Configuration File

1. Click on the "Choose File" button in the "Upload YAML Configuration File" section.
2. Select your YAML configuration file.
3. Click the "Generate Config" button.

### Option 2: Enter Configuration Manually

1. Fill in the following fields:
   - Catchall Port: The default port for your NGINX server (e.g., 7000)
   - App Name: Name of your application
   - App FQDN: Fully Qualified Domain Names (comma-separated for multiple domains)
   - App Runtime Port: The port your application runs on

2. Add Path-Based Access Restrictions:
   - Click "Add Path Restriction"
   - Enter the path (e.g., "/api")
   - Enter the IP filter name to apply to this path

3. Add IP Filters:
   - Click "Add IP Filter"
   - Enter a name for the filter
   - Enter IP addresses or ranges (comma-separated)

4. Click the "Generate Config" button.

### View Generated Configuration

After clicking "Generate Config", the tool will display the generated NGINX configuration below the form. You can copy this configuration and use it in your NGINX setup.

## YAML File Structure

If you're uploading a YAML file, it should follow this structure:

```yaml
catchall:
  default:
    port: 7000

app:
  yourappname:
    catchall: "default"
    fqdn:
      - "yourdomain.com"
      - "www.yourdomain.com"
    runtime_port: 8000
    path_based_access_restriction:
      "/": 
        ipFilter: "allowall"
      "/admin":
        ipFilter: "adminonly"

ipfilter:
  allowall:
    - "0.0.0.0/0"
    - "::/0"
  adminonly:
    - "192.168.1.0/24"
    - "2001:db8::/32"
```

## License

This project is licensed under the MIT License.