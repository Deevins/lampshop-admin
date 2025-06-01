// src/pages/Auth/Login.tsx
import React, {useState} from "react";
import {useNavigate} from "react-router-dom";
import styles from "./Login.module.scss";
import {login} from "../../api/AuthApi.ts";


const Login: React.FC = () => {
    const navigate = useNavigate();
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            const token = await login(username, password);
            localStorage.setItem("jwtToken", token);
            navigate("/products");
        } catch (err: any) {
            console.error("Login error:", err);
            setError("Неверный логин или пароль");
            setLoading(false);
        }
    };

    return (
        <div className={styles["login-container"]}>
            <h2 className={styles.title}>Войти как администратор</h2>
            <form onSubmit={handleSubmit} className={styles.form}>
                <div className={styles["form-group"]}>
                    <label htmlFor="username">Логин</label>
                    <input
                        id="username"
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                    />
                </div>
                <div className={styles["form-group"]}>
                    <label htmlFor="password">Пароль</label>
                    <input
                        id="password"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                {error && <p className={styles.error}>{error}</p>}
                <button
                    type="submit"
                    disabled={loading}
                    className={styles["submit-btn"]}
                >
                    {loading ? "Загрузка..." : "Войти"}
                </button>
            </form>
        </div>
    );
};

export default Login;
