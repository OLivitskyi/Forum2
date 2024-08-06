// Function to display an error message
export const showError = (message) => {
    const errorMessage = document.getElementById('error-message');
    if (errorMessage) {
        errorMessage.textContent = message;
    } else {
        console.error("Error element not found");
    }
};

// Function to clear the error message
export const clearError = () => {
    const errorMessage = document.getElementById('error-message');
    if (errorMessage) {
        errorMessage.textContent = '';
    } else {
        console.error("Error element not found");
    }
};
