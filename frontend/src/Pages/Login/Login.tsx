import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import styles from "./Login.module.scss";

const Login: React.FC = () => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState<string | null>(null);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        const success = await login(username, password);
        if (success) {
            navigate("/products");
        } else {
            setError("Неверный логин или пароль");
        }
    };

    return (
        <div className={styles["login-container"]}>
            <div className={styles["login-box"]}>
                <h2 className={styles["login-title"]}>Вход</h2>
                <form onSubmit={handleSubmit}>
                    <div className={styles["input-group"]}>
                        <label className={styles["input-label"]} htmlFor="username">
                            Логин
                        </label>
                        <input
                            id="username"
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            className={styles["input-field"]}
                            required
                        />
                    </div>
                    <div className={styles["input-group"]}>
                        <label className={styles["input-label"]} htmlFor="password">
                            Пароль
                        </label>
                        <input
                            id="password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className={styles["input-field"]}
                            required
                        />
                    </div>
                    {error && <div className={styles["error-message"]}>{error}</div>}
                    <button type="submit" className={styles["login-button"]}>
                        Войти
                    </button>
                </form>
            </div>
        </div>
    );
};

export default Login;
