// TODO :: create skeleton.tsx & error-boundary.tsx

type ProjectListingProps = Readonly<{
    dummy?: string;
}>;

export default function ProjectListing(props: ProjectListingProps) {
    console.info('ProjectListing.props :: ', props);

    return (
        <>
            <h1>Project Listing</h1>
        </>
    );
}
