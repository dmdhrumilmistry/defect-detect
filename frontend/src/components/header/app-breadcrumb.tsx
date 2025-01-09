import type { LoaderReturnValue, Maybe, RouteHandle } from '@/types';
import { Fragment } from 'react';
import { Link, type UIMatch, useMatches } from 'react-router-dom';
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';

export default function AppBreadcrumb() {
    const matches = useMatches() as UIMatch<Maybe<LoaderReturnValue>, RouteHandle>[];
    const breadcrumbs = matches
        .filter((match) => Boolean(match.handle?.breadcrumb))
        .map((match) => match.handle.breadcrumb?.(match.data));
    const lastBreadcrumb = breadcrumbs.pop();

    return (
        lastBreadcrumb && (
            <Breadcrumb>
                <BreadcrumbList>
                    {breadcrumbs.map(
                        (breadcrumb) =>
                            breadcrumb && (
                                <Fragment key={breadcrumb.label}>
                                    <BreadcrumbItem className="hidden md:block">
                                        <BreadcrumbLink asChild>
                                            <Link to={breadcrumb.href}>{breadcrumb?.label}</Link>
                                        </BreadcrumbLink>
                                    </BreadcrumbItem>
                                    <BreadcrumbSeparator className="hidden md:block" />
                                </Fragment>
                            )
                    )}
                    <BreadcrumbItem>
                        <BreadcrumbPage>{lastBreadcrumb.label}</BreadcrumbPage>
                    </BreadcrumbItem>
                </BreadcrumbList>
            </Breadcrumb>
        )
    );
}
