const ipFiltersContainer = document.getElementById('ip_filters_container');
const addIpFilterBtn = document.getElementById('add_ip_filter');
const pathRestrictionsContainer = document.getElementById('path_restrictions_container');
const addPathRestrictionBtn = document.getElementById('add_path_restriction');
const configForm = document.getElementById('configForm');
const resultContainer = document.getElementById('result');
const nginxConfigPre = document.getElementById('nginx_config');
const yamlFileInput = document.getElementById('yaml_file');
const manualInputs = document.querySelectorAll('#catchall_port, #app_name, #app_fqdn, #app_port');

let ipFilterCount = 0;
let pathRestrictionCount = 0;

function addIpFilter() {
    const filterId = `ip_filter_${ipFilterCount}`;
    const filterHtml = `
        <div id="${filterId}" class="mb-2 p-2 border rounded">
            <div class="flex items-center mb-2">
                <input class="flex-grow mr-2 shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" placeholder="Filter Name" data-filter-name>
                <input class="flex-grow mr-2 shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" placeholder="IP Addresses (comma-separated)" data-filter-ips>
                <button type="button" class="bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 rounded focus:outline-none focus:shadow-outline" onclick="removeIpFilter('${filterId}')">Remove</button>
            </div>
        </div>
    `;
    ipFiltersContainer.insertAdjacentHTML('beforeend', filterHtml);
    ipFilterCount++;
    clearFileInput();
}

function removeIpFilter(filterId) {
    document.getElementById(filterId).remove();
}

function addPathRestriction() {
    const restrictionId = `path_restriction_${pathRestrictionCount}`;
    const restrictionHtml = `
        <div id="${restrictionId}" class="mb-2 p-2 border rounded">
            <div class="flex items-center mb-2">
                <input class="flex-grow mr-2 shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" placeholder="Path (e.g., /)" data-path>
                <input class="flex-grow mr-2 shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" placeholder="IP Filter Name" data-ip-filter>
                <button type="button" class="bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 rounded focus:outline-none focus:shadow-outline" onclick="removePathRestriction('${restrictionId}')">Remove</button>
            </div>
        </div>
    `;
    pathRestrictionsContainer.insertAdjacentHTML('beforeend', restrictionHtml);
    pathRestrictionCount++;
    clearFileInput();
}

function removePathRestriction(restrictionId) {
    document.getElementById(restrictionId).remove();
}

function clearFileInput() {
    yamlFileInput.value = '';
}

function clearManualInputs() {
    manualInputs.forEach(input => input.value = '');
    while (ipFiltersContainer.firstChild) {
        ipFiltersContainer.removeChild(ipFiltersContainer.firstChild);
    }
    while (pathRestrictionsContainer.firstChild) {
        pathRestrictionsContainer.removeChild(pathRestrictionsContainer.firstChild);
    }
    ipFilterCount = 0;
    pathRestrictionCount = 0;
}

yamlFileInput.addEventListener('change', (e) => {
    if (e.target.files.length > 0) {
        clearManualInputs();
    }
});

manualInputs.forEach(input => {
    input.addEventListener('input', clearFileInput);
});

configForm.addEventListener('submit', async function(e) {
    e.preventDefault();
    const submitButton = this.querySelector('button[type="submit"]');

    // Show loading state
    submitButton.disabled = true;
    submitButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Generating...';

    const formData = new FormData();
    const yamlFile = yamlFileInput.files[0];
    
    if (yamlFile) {
        formData.append('yaml_file', yamlFile);
    } else {
        const config = {
            catchall: {
                default: {
                    port: parseInt(document.getElementById('catchall_port').value) || 7000
                }
            },
            app: {
                [document.getElementById('app_name').value]: {
                    catchall: "default",
                    fqdn: document.getElementById('app_fqdn').value.split(',').map(s => s.trim()),
                    runtime_port: parseInt(document.getElementById('app_port').value) || 8000,
                    path_based_access_restriction: {}
                }
            },
            ipfilter: {}
        };

        // Collect IP filters
        document.querySelectorAll('[id^="ip_filter_"]').forEach(filter => {
            const name = filter.querySelector('[data-filter-name]').value;
            const ips = filter.querySelector('[data-filter-ips]').value.split(',').map(ip => ip.trim());
            if (name && ips.length > 0) {
                config.ipfilter[name] = ips;
            }
        });

        // Collect path-based access restrictions
        document.querySelectorAll('[id^="path_restriction_"]').forEach(restriction => {
            const path = restriction.querySelector('[data-path]').value;
            const ipFilter = restriction.querySelector('[data-ip-filter]').value;
            if (path && ipFilter) {
                config.app[document.getElementById('app_name').value].path_based_access_restriction[path] = {
                    ipFilter: ipFilter
                };
            }
        });

        formData.append('config', JSON.stringify(config));
    }

    try {
        const response = await fetch('/generate', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        nginxConfigPre.textContent = data.config;
        resultContainer.classList.remove('hidden');
    } catch (error) {
        console.error('Error:', error);
        nginxConfigPre.textContent = 'Error generating configuration: ' + error.message;
    } finally {
        // Reset button state
        submitButton.disabled = false;
        submitButton.textContent = 'Generate Config';
    }
});

// Initialize with one IP filter and one path restriction
addIpFilter();
addPathRestriction();

// Add event listeners for the "Add" buttons
addIpFilterBtn.addEventListener('click', addIpFilter);
addPathRestrictionBtn.addEventListener('click', addPathRestriction);