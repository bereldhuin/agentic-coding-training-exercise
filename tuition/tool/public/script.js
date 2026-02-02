document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('register-form');
    const serverList = document.getElementById('server-list');

    // Fetch and display servers
    const fetchServers = async () => {
        try {
            const response = await fetch('/api/servers');
            const servers = await response.json();
            renderServers(servers);
        } catch (error) {
            console.error('Error fetching servers:', error);
        }
    };

    const renderServers = (servers) => {
        serverList.innerHTML = servers.map(server => `
            <li class="${server.status}">
                <div class="server-info">
                    <strong>${server.name}</strong><br>
                    <small>${server.url}</small><br>
                    <span class="status-badge status-${server.status}">${server.status}</span>
                    ${server.lastChecked ? `<small>Last checked: ${new Date(server.lastChecked).toLocaleString()}</small>` : ''}
                </div>
                <div class="server-actions">
                    <button onclick="verifyServer('${server.id}')">Verify</button>
                    <button onclick="copyMcpCommand('${server.name}', '${server.url}')">Copy Command</button>
                </div>
            </li>
        `).join('');
    };

    // Register new server
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('name').value;
        const url = document.getElementById('url').value;

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, url })
            });
            if (response.ok) {
                registerForm.reset();
                fetchServers();
            }
        } catch (error) {
            console.error('Error registering server:', error);
        }
    });

    // Verify server
    window.verifyServer = async (id) => {
        try {
            const response = await fetch(`/api/verify/${id}`, { method: 'POST' });
            if (response.ok) {
                fetchServers();
            }
        } catch (error) {
            console.error('Error verifying server:', error);
        }
    };

    // Copy MCP command to clipboard
    window.copyMcpCommand = async (name, url) => {
        try {
            // Extract IP and port from URL (e.g., http://192.168.1.100:3000/mcp -> 192.168.1.100:3000)
            const urlObj = new URL(url);
            const ipWithPort = `${urlObj.hostname}:${urlObj.port}`;

            // Build the command
            const command = `claude mcp add ${name} ${ipWithPort} --transport http --scope project`;

            // Copy to clipboard
            await navigator.clipboard.writeText(command);

            // Show success feedback
            alert('Command copied to clipboard!');
        } catch (error) {
            console.error('Error copying to clipboard:', error);
            alert('Failed to copy command');
        }
    };

    // Copy text to clipboard (generic helper)
    window.copyToClipboard = async (text) => {
        try {
            await navigator.clipboard.writeText(text);
            alert('Copied to clipboard!');
        } catch (error) {
            console.error('Error copying to clipboard:', error);
            alert('Failed to copy');
        }
    };

    // Initial load
    fetchServers();
});
