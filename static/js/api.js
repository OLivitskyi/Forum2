export const getCategories = async () => {
    const cookieString = document.cookie;
    const cookies = cookieString.split('; ').reduce((acc, cookie) => {
        const [name, value] = cookie.split('=');
        acc[name] = value;
        return acc;
    }, {});
    
    const sessionToken = cookies['session_token'];

    if (!sessionToken) {
        throw new Error("No session token found");
    }

    const response = await fetch("/api/get-categories", {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${sessionToken}`
        }
    });
    if (!response.ok) {
        throw new Error("Failed to fetch categories");
    }
    return await response.json();
};
