# Security Policy

## Supported Versions

The following table shows which versions of the project are currently supported with security updates.

| Version | Supported          |
| -------- | ------------------ |
| 1.x.x    | :white_check_mark: |
| 0.x.x    | :x:                |

Only the latest **stable version** receives security updates and vulnerability patches.  
Older or experimental branches may not receive fixes unless explicitly stated.

---

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please **do not open a public issue**.  
Instead, follow the process below to ensure responsible disclosure.

### 🔒 How to Report

- Send an email to **security@yourdomain.com** (or your GitHub username if you don’t have a domain yet).  
- Include as much detail as possible:
  - Steps to reproduce the issue  
  - Potential impact (e.g., RCE, DoS, data leak, etc.)  
  - A proof of concept (PoC) if available  

### 🕒 Response Expectations

- You will receive an initial acknowledgment within **72 hours**.  
- A member of the maintainers’ team will:
  - Confirm the issue,
  - Assess its severity,
  - Work on a fix or mitigation plan.

### 🔧 Disclosure Policy

Once a vulnerability has been fixed:
- A new release will be published with a security patch note (e.g., `v1.2.3-security.1`).
- You will be credited in the changelog (if desired).
- We’ll publish the details **after** users have had sufficient time to update.

---

## Security Best Practices

To help secure your own deployments:
- Always use HTTPS and strong authentication.
- Keep your Go environment and dependencies up to date.
- Restrict access to admin and monitoring endpoints.
- Use environment variables for secrets (never commit credentials).

---

## Contact

If you have any questions about this policy, contact:
- GitHub: [@your-github-username](https://github.com/your-github-username)
- Email: security@yourdomain.com
