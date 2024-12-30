// TODO :: create skeleton.tsx & error-boundary.tsx

type ProjectProps = Readonly<{
    dummy: string;
}>;

export default function Project(props: ProjectProps) {
    console.info('Project.props :: ', props);

    return (
        <>
            <h1>Dummy Project</h1>
        </>
    );
}
