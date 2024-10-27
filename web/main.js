document.addEventListener("DOMContentLoaded", function() {
    
    function startMonitoring() {
      fetch('/start-monitoring')
        .then(response => {
          if (response.ok) {
            alert("Monitoring started successfully!");
            fetchLogs(); // Carrega os logs automaticamente após o monitoramento
          } else {
            alert("Failed to start monitoring.");
          }
        })
        .catch(error => {
          console.error("Error:", error);
          alert("Failed to connect to the server.");
        });
    }
    function listLogs() {
      fetch('/list-logs')
        .then(response => response.text())
        .then(data => {
          const logsDiv = document.getElementById("logs");
          logsDiv.innerHTML = ""; // Limpa logs
  
          const logFiles = data.trim().split("\n");
          logFiles.forEach(file => {
            if (file) {
              const logButton = document.createElement("button");
              logButton.textContent = file;
              logButton.onclick = () => viewLog(file);
              logsDiv.appendChild(logButton);
            }
          });
        })
        .catch(error => {
          console.error("Error listing logs:", error);
        });
    }
  
    
    function viewLog(file) {
      fetch(`/view-log?file=${file}`)
        .then(response => response.text())
        .then(data => {
          const logsDiv = document.getElementById("logs");
          logsDiv.innerHTML = `<h3>Contents of ${file}:</h3><pre>${data}</pre>`;
        })
        .catch(error => {
          console.error("Error viewing log:", error);
        });
    }
  
    
    function fetchLogs() {
      const logsDiv = document.getElementById("logs");
      logsDiv.innerHTML = "<p>Loading logs...</p>"; 
  
      fetch('/logs')
        .then(response => response.text())
        .then(data => {
          logsDiv.innerHTML = "";
  
          if (data.trim() === "") {
            logsDiv.innerHTML = "<p>No logs available.</p>"; // Exibe uma mensagem se não houver logs
            return;
          }
  
          const logsArray = data.split("\n"); 
          logsArray.forEach(log => {
            if (log.trim()) { 
              const logElement = document.createElement("p");
              logElement.textContent = log;
              logsDiv.appendChild(logElement);
            }
          });
        })
        .catch(error => {
          console.error("Error:", error);
          logsDiv.innerHTML = "<p>Failed to fetch logs. Please try again later.</p>";
        });
    }
  
    
    function exitProgram() {
      fetch('/exit')
        .then(response => {
          if (response.ok) {
            alert("Program exited successfully!");
          } else {
            alert("Failed to exit the program.");
          }
        })
        .catch(error => {
          console.error("Error:", error);
          alert("Failed to connect to the server.");
        });
    }
  
    
    document.getElementById("start-monitoring").addEventListener("click", startMonitoring);
    document.getElementById("list-logs").addEventListener("click", listLogs);
    document.getElementById("exit-program").addEventListener("click", exitProgram);
  });
  