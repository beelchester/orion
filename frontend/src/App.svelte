<script lang="ts">
  import { 
    Download, 
    GetDownloadProgress, 
    PauseDownload, 
    ResumeDownload, 
    CancelDownload 
  } from '../wailsjs/go/main/App.js';
  
  let resultText = "Enter the download link";
  let url = "";
  let progress = 0;
  let downloading = false;
  let paused = false;
  let gid = "";
  let downloadSpeed = 0;
  let totalSize = 0;
  let completedSize = 0;
  let status = "";
  let showProgress = false;  // New flag to control progress bar visibility
  
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
    resultText = "Starting download...";
    showProgress = true;  // Show progress bar when download starts
    
    try {
      // Start download and get GID
      gid = await Download(url);
      resultText = "Download started";
      trackProgress();
    } catch (error) {
      resultText = `Error: ${error.message || error}`;
      downloading = false;
      showProgress = false;  // Hide progress bar on error
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
      showProgress = false;  // Hide progress bar when download is canceled
      resultText = "Download canceled!";
      
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
          
          // Update paused status based on returned status
          paused = status === "paused";
          
          resultText = paused 
            ? `Download paused: ${formatBytes(completedSize)} of ${formatBytes(totalSize)}`
            : `Downloading: ${formatBytes(completedSize)} of ${formatBytes(totalSize)} (${formatSpeed(downloadSpeed)})`;
          
          if (progress >= 100 || status === "complete") {
            downloading = false;
            paused = false;
            resultText = `Download Complete! Total size: ${formatBytes(totalSize)}`;
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
  <div class="result">{resultText}</div>
  
  <div class="input-box">
    <input
      bind:value={url}
      class="input"
      type="text"
      placeholder="Enter URL to download..."
      on:keypress={(e) => e.key === 'Enter' && !downloading && download()}
    />
    <button class="btn" on:click={download} disabled={downloading}>
      {downloading ? "Downloading..." : "Download"}
    </button>
  </div>
  
  {#if showProgress}
    <div class="progress-container">
      <div class="progress-stats">
        <div class="stat-item">
          <span class="stat-label">Size:</span>
          <span class="stat-value">{formatBytes(completedSize)} / {formatBytes(totalSize)}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Speed:</span>
          <span class="stat-value">{formatSpeed(downloadSpeed)}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Status:</span>
          <span class="stat-value">{status}</span>
        </div>
      </div>
      
      <div class="progress-bar">
        <div class="progress" style="width: {progress}%" class:paused={paused}>
          <span class="progress-text">{Math.round(progress)}%</span>
        </div>
      </div>
      
      <!-- Control buttons -->
      <div class="control-buttons">
        {#if paused}
          <button class="control-btn resume" on:click={resume}>
            Resume
          </button>
        {:else if downloading}
          <button class="control-btn pause" on:click={pause}>
            Pause
          </button>
        {/if}
        
        {#if downloading || paused}
          <button class="control-btn cancel" on:click={cancel}>
            Cancel
          </button>
        {/if}
      </div>
    </div>
  {/if}
</main>
 
<style>
  main {
    padding: 20px;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .result {
    margin-bottom: 20px;
    padding: 10px;
    background: #f5f5f5;
    border-radius: 5px;
    min-height: 20px;
    color: black;
  }
  
  .input-box {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }
  
  .input {
    flex: 1;
    padding: 8px 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  .btn {
    padding: 8px 16px;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .btn:disabled {
    background: #cccccc;
    cursor: not-allowed;
  }
  
  .progress-container {
    margin-top: 20px;
    border: 1px solid #ddd;
    border-radius: 6px;
    padding: 12px;
    background-color: #f9f9f9;
  }
  
  .progress-stats {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .stat-item {
    display: flex;
    flex-direction: column;
    min-width: 100px;
  }
  
  .stat-label {
    font-size: 12px;
    color: #666;
    margin-bottom: 2px;
  }
  
  .stat-value {
    font-weight: bold;
    color: #333;
  }
  
  .progress-bar {
    width: 100%;
    height: 24px;
    background: #f0f0f0;
    border-radius: 5px;
    position: relative;
    overflow: hidden;
    border: 1px solid #ddd;
    margin-bottom: 15px;
  }
  
  .progress {
    height: 100%;
    background: #4caf50;
    border-radius: 4px;
    transition: width 0.5s ease-in-out;
    min-width: 30px;
  }
  
  .progress.paused {
    background: #ffc107;  /* Yellow color when paused */
  }
  
  .progress-text {
    position: absolute;
    width: 100%;
    text-align: center;
    font-size: 14px;
    color: #333;
    font-weight: bold;
    line-height: 24px;
  }
  
  .control-buttons {
    display: flex;
    gap: 10px;
    margin-top: 10px;
  }
  
  .control-btn {
    padding: 6px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    min-width: 80px;
  }
  
  .pause {
    background: #ffc107;
    color: #333;
  }
  
  .resume {
    background: #2196f3;
    color: white;
  }
  
  .cancel {
    background: #f44336;
    color: white;
  }
</style>