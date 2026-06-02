import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { NotificationProvider } from './contexts/NotificationContext';
import { VersionProvider } from './contexts/VersionContext';
import { useAuth } from './contexts/useAuth';
import AuthLayout from './layouts/AuthLayout';
import DashboardLayout from './layouts/DashboardLayout';
import LoginPage from './pages/Auth/LoginPage';
import RegisterPage from './pages/Auth/RegisterPage';
import OAuthCallbackPage from './pages/Auth/OAuthCallbackPage';
import DashboardPage from './pages/Dashboard/DashboardPage';
import ChangelogPage from './pages/Changelog/ChangelogPage';
import ProjectsPage from './pages/Projects/ProjectsPage';
import ProjectDetailPage from './pages/Projects/ProjectDetailPage';
import TaskDetailPage from './pages/Projects/TaskDetailPage';
import MyTasksPage from './pages/Tasks/MyTasksPage';
import ProfilePage from './pages/Profile/ProfilePage';
import ProtectedRoute from './components/Auth/ProtectedRoute';

function AuthRedirect({ children }: { children: React.ReactNode }) {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (isAuthenticated) {
        return <Navigate to="/dashboard" replace />;
    }

    return <>{children}</>;
}

function App() {
    return (
        <BrowserRouter>
            <AuthProvider>
                <NotificationProvider>
                <Routes>
                    <Route path="/login" element={
                        <AuthRedirect><AuthLayout><LoginPage /></AuthLayout></AuthRedirect>
                    } />
                    <Route path="/register" element={
                        <AuthRedirect><AuthLayout><RegisterPage /></AuthLayout></AuthRedirect>
                    } />
                    <Route path="auth/google/callback" element={<OAuthCallbackPage />} />
                    <Route path="auth/yandex/callback" element={<OAuthCallbackPage />} />
                    <Route path="/" element={
                        <ProtectedRoute><VersionProvider><DashboardLayout /></VersionProvider></ProtectedRoute>
                    }>
                        <Route index element={<Navigate to="/dashboard" replace />} />
                        <Route path="dashboard" element={<DashboardPage />} />
                        <Route path="changelog" element={<ChangelogPage />} />
                        <Route path="projects" element={<ProjectsPage />} />
                        <Route path="projects/:id" element={<ProjectDetailPage />} />
                        <Route path="projects/:id/tasks/:taskId" element={<TaskDetailPage />} />
                        <Route path="tasks" element={<MyTasksPage />} />
                        <Route path="profile" element={<ProfilePage />} />
                    </Route>
                    <Route path="*" element={<Navigate to="/dashboard" replace />} />
                </Routes>
                </NotificationProvider>
            </AuthProvider>
        </BrowserRouter>
    );
}

export default App;
