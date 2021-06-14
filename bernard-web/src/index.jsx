import Router from 'preact-router';
import { h, render } from 'preact';

const Hosts = () => {
    return <div>hosts</div>
}

const Checks = () => {
    return <div>Checks</div>
}

const Main = () => (
    <Router>
        <Hosts path="/" />
        <Checks path="/checks" />
    </Router>
);

render(<Main />, document.body);