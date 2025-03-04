<script lang="ts">
  import { 
    Download, 
    GetDownloadProgress, 
    PauseDownload, 
    ResumeDownload, 
    CancelDownload 
  } from '../wailsjs/go/main/App.js';
  import './App.css';
  
  let resultText = "Ready to download";
  let url = "";
  let progress = 0;
  let downloading = false;
  let paused = false;
  let gid = "";
  let downloadSpeed = 0;
  let totalSize = 0;
  let completedSize = 0;
  let status = "";
  let showProgress = false;
  let estimatedTimeLeft = "calculating...";
  let startTime = 0;
  
  // converting bytes into human-readable format
  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }
  
  // format speed into human-readable format
  function formatSpeed(bytesPerSecond) {
    return formatBytes(bytesPerSecond) + '/s';
  }
  
  // estimated time remaining/arival
  function calculateETA(totalSize, completedSize, speed) {
    if (speed <= 0 || totalSize <= 0) return "calculating...";
    
    const remainingBytes = totalSize - completedSize;
    const secondsLeft = remainingBytes / speed;
    
    if (secondsLeft < 60) {
      return `${Math.round(secondsLeft)} seconds`;
    } else if (secondsLeft < 3600) {
      return `${Math.round(secondsLeft / 60)} minutes`;
    } else {
      const hours = Math.floor(secondsLeft / 3600);
      const minutes = Math.round((secondsLeft % 3600) / 60);
      return `${hours}h ${minutes}m`;
    }
  }
  
  // Get filename from URL
  function getFilenameFromUrl(url) {
    try {
      const urlObj = new URL(url);
      const pathname = urlObj.pathname;
      const filename = pathname.split('/').pop();
      return filename || "file";
    } catch (e) {
      return "file";
    }
  }
  
  async function download() {
    if (!url) {
      resultText = "Please enter a download link";
      return;
    }
    
    downloading = true;
    paused = false;
    progress = 0;
    downloadSpeed = 0;
    totalSize = 0;
    completedSize = 0;
    status = "Starting";
    resultText = `Initializing download for ${getFilenameFromUrl(url)}...`;
    showProgress = true;
    startTime = Date.now();
    
    try {
      // Start download and get GID
      gid = await Download(url);
      resultText = "Download in progress";
      trackProgress();
    } catch (error) {
      resultText = `Error: ${error.message || error}`;
      downloading = false;
      showProgress = false;
    }
  }
  
  async function pause() {
    if (!gid || !downloading) return;
    
    try {
      await PauseDownload(gid);
      paused = true;
      resultText = "Download paused";
    } catch (error) {
      resultText = `Error pausing: ${error.message || error}`;
    }
  }
  
  async function resume() {
    if (!gid || !paused) return;
    
    try {
      await ResumeDownload(gid);
      paused = false;
      resultText = "Download resumed";
    } catch (error) {
      resultText = `Error resuming: ${error.message || error}`;
    }
  }
  
  async function cancel() {
    if (!gid) return;
    
    try {
      await CancelDownload(gid);
      downloading = false;
      paused = false;
      showProgress = false;
      resultText = "Download canceled";
      
      // Reset download-related variables
      progress = 0;
      downloadSpeed = 0;
      completedSize = 0;
      totalSize = 0;
      status = "";
      gid = "";
    } catch (error) {
      resultText = `Error canceling: ${error.message || error}`;
    }
  }
  
  async function trackProgress() {
    while (downloading) {
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      try {
        if (!gid) {
          continue;
        }
        
        let info = await GetDownloadProgress(gid);
        
        if (info) {
          progress = Math.min(100, Math.max(0, info.progress));
          downloadSpeed = info.downloadSpeed;
          totalSize = info.totalSize;
          completedSize = info.completedSize;
          status = info.status;
          
          // Update estimated time
          estimatedTimeLeft = calculateETA(totalSize, completedSize, downloadSpeed);
          
          // Update paused status based on returned status
          paused = status === "paused";
          
          resultText = paused 
            ? `Download paused: ${formatBytes(completedSize)} of ${formatBytes(totalSize)}`
            : `Downloading: ${getFilenameFromUrl(url)}`;
          
          if (progress >= 100 || status === "complete") {
            downloading = false;
            paused = false;
            const totalTime = ((Date.now() - startTime) / 1000).toFixed(1);
            resultText = `Download Complete! Total size: ${formatBytes(totalSize)} in ${totalTime}s`;
          }
        }
      } catch (error) {
        if (error.toString().includes("not found") || error.toString().includes("invalid GID")) {
          downloading = false;
          paused = false;
          resultText = "Download completed or GID invalid.";
          progress = 100; 
        }
      }
    }
  }
</script>
 
<main>
  <div class="app-container">
    <h1 class="app-title">Orion</h1>
    
    <div class="input-container">
      <div class="url-input-wrapper">
        <input
          bind:value={url}
          class="url-input"
          type="text"
          placeholder="Enter URL to download..."
          on:keypress={(e) => e.key === 'Enter' && !downloading && download()}
        />
        <button class="download-btn" on:click={!downloading ? download : paused ? resume : pause} disabled={downloading && !paused && progress >= 100}>
          {#if !downloading}
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
            Download
          {:else if paused}
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
            Resume
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="6" y="4" width="4" height="16"></rect><rect x="14" y="4" width="4" height="16"></rect></svg>
            Pause
          {/if}
        </button>
      </div>
    </div>
    
    <div class="status-container" class:active={resultText !== "Ready to download"}>
      <div class="status-indicator" class:success={status === "complete"} class:warning={paused} class:error={resultText.includes("Error")}>
        {#if status === "complete"}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
        {:else if paused}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="6" y="4" width="4" height="16"></rect><rect x="14" y="4" width="4" height="16"></rect></svg>
        {:else if resultText.includes("Error")}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
        {:else if downloading}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"></path></svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
        {/if}
      </div>
      <div class="status-text">{resultText}</div>
    </div>
    
    {#if showProgress}
      <div class="progress-panel">
        <div class="progress-container">
          <div class="progress-bar">
            <div class="progress-fill" style="width: {progress}%" class:paused={paused} class:completed={status === "complete"}></div>
            <span class="progress-text">{Math.round(progress)}%</span>
          </div>
        </div>
        
        <div class="progress-stats">
          <div class="stat-card">
            <div class="stat-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>
            </div>
            <div class="stat-details">
              <div class="stat-label">Speed</div>
              <div class="stat-value">{formatSpeed(downloadSpeed)}</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            </div>
            <div class="stat-details">
              <div class="stat-label">ETA</div>
              <div class="stat-value">{estimatedTimeLeft}</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path></svg>
            </div>
            <div class="stat-details">
              <div class="stat-label">Size</div>
              <div class="stat-value">{formatBytes(completedSize)} / {formatBytes(totalSize)}</div>
            </div>
          </div>
        </div>
        
        <div class="control-buttons">
          {#if downloading || paused}
            <button class="control-btn cancel" on:click={cancel}>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
              Cancel
            </button>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</main>
 
