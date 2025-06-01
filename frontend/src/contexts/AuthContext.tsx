import React, {createContext, useContext, useState} from "react";
import {useNavigate} from "react-router-dom";

interface AuthContextType {
    user: string | null;
    login: (username: string, password: string) => Promise<boolean>;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
                                                                          children
                                                                      }) => {
    const [user, setUser] = useState<string | null>(
        localStorage.getItem("admin_user")
    );
    const navigate = useNavigate();

    const login = async (username: string, password: string) => {
        if (username === "admin" && password === "password123") {
            localStorage.setItem("admin_user", username);
            setUser(username);
            return true;
        }
        return false;
    };

    const logout = () => {
        localStorage.removeItem("admin_user");
        setUser(null);
        navigate("/login");
    };

    return (
        <AuthContext.Provider value={{user, login, logout}}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const ctx = useContext(AuthContext);
    if (!ctx) {
        throw new Error("useAuth must be used within AuthProvider");
    }
    return ctx;
};
