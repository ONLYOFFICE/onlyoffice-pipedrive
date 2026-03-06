# ONLYOFFICE app for Pipedrive

Work on deal documents without leaving Pipedrive. This integration adds [ONLYOFFICE Docs](https://www.onlyoffice.com/docs) to your [Pipedrive](https://www.pipedrive.com/) deals so you can create, edit, and co-author office documents — all saved back to Pipedrive in real time.

<!-- Banner -->
<p align="center">
  <a href="https://www.onlyoffice.com/office-for-pipedrive">
    <img src="https://static-blog.onlyoffice.com/wp-content/uploads/2023/08/15122902/ONLYOFFICE-Pipedrive-connector.jpg" alt="ONLYOFFICE for Pipedrive" />
  </a>
</p>

## Top features ✨

- **Edit inside deals**: Open docs from the ONLYOFFICE Documents section in a deal — no switching tabs.
- **Create or upload**: Add new documents, spreadsheets, and presentations or upload existing files to the deal.
- **Real-time collaboration**: Co-edit together, comment, and chat directly in the editors.
- **Secure by default**: JWT protection is enabled in ONLYOFFICE Docs.

## Installing ONLYOFFICE Docs

To be able to work with office files within Pipedrive, you will need an instance of ONLYOFFICE Docs. You can install the self-hosted version of the editors or opt for ONLYOFFICE Docs Cloud which doesn't require downloading and installation.

### Self-hosted editors

You can install free Community version of ONLYOFFICE Docs or scalable Enterprise Edition.

To install free Community version, use [Docker](https://github.com/onlyoffice/Docker-DocumentServer) (recommended) or follow [these instructions](https://helpcenter.onlyoffice.com/docs/installation/docs-community-install-ubuntu.aspx) for Debian, Ubuntu, or derivatives.

To install Enterprise Edition, follow the instructions [here](https://helpcenter.onlyoffice.com/docs/installation/enterprise).

Community Edition vs Enterprise Edition comparison can be found [here](#onlyoffice-docs-editions).

### ONLYOFFICE Docs Cloud

To get ONLYOFFICE Docs Cloud, get started [here](https://www.onlyoffice.com/docs-registration).

## App installation and configuration ⚙️

You can add the ONLYOFFICE app from the [Pipedrive App Marketplace](https://www.pipedrive.com/en/marketplace/app/onlyoffice/ab1362af3cc9f7be). 

Once done, go to the ONLYOFFICE app settings page (**Tools and Integrations -> Marketplace apps -> ONLYOFFICE**) and enter the name of the server with ONLYOFFICE Docs installed in the *Document Server Address* field.

Starting from version 7.2 of ONLYOFFICE Docs, JWT is enabled by default and the secret key is generated automatically to restrict the access to the editors and for security reasons and data integrity. You can specify your own *Document Server Secret* on the settings page. In the ONLYOFFICE Docs [config file](https://api.onlyoffice.com/docs/docs-api/additional-api/signature/), specify the same secret key to enable the validation.

## App usage 

The app allows working with office documents directly within the Pipedrive frontend.

You can create and upload text documents, spreadsheets, and presentations within your Pipedrive Deals. Just click the corresponding button (**Create or upload document**) in the ONLYOFFICE Documents section.

To edit the created files, reach to the **ONLYOFFICE Documents** section and open the needed document by clicking the pencil icon. Everyone who has access to the deal can open the file for editing. You can also collaborate on documents in real time together with your colleagues.

## ONLYOFFICE Docs editions

Self-hosted **ONLYOFFICE Docs** is packaged as Document Server:

* Community Edition 🆓 (`onlyoffice-documentserver` package)
* Enterprise Edition 🏢 (`onlyoffice-documentserver-ee` package)

The table below will help you to make the right choice.

| Pricing and licensing | Community Edition | Enterprise Edition |
| ------------- | ------------- | ------------- |
| | [Get it now](https://www.onlyoffice.com/download-community?utm_source=github&utm_medium=cpc&utm_campaign=GitHubPipedrive#docs-community)  | [Start Free Trial](https://www.onlyoffice.com/download?utm_source=github&utm_medium=cpc&utm_campaign=GitHubPipedrive#docs-enterprise)  |
| Cost  | FREE  | [Go to the pricing page](https://www.onlyoffice.com/docs-enterprise-prices?utm_source=github&utm_medium=cpc&utm_campaign=GitHubPipedrive)  |
| Simultaneous connections | up to 20 maximum  | As in chosen pricing plan |
| Number of users | up to 20 recommended | As in chosen pricing plan |
| License | GNU AGPL v.3 | Proprietary |
| **Support** | **Community Edition** | **Enterprise Edition** |
| Documentation | [Help Center](https://helpcenter.onlyoffice.com/docs/installation/community) | [Help Center](https://helpcenter.onlyoffice.com/docs/installation/enterprise) |
| Standard support | [GitHub](https://github.com/ONLYOFFICE/DocumentServer/issues) or paid | 1 or 3 years support included |
| Premium support | [Contact us](mailto:sales@onlyoffice.com) | [Contact us](mailto:sales@onlyoffice.com) |
| **Services** | **Community Edition** | **Enterprise Edition** |
| Conversion Service                | + | + |
| Document Builder Service          | + | + |
| **Interface** | **Community Edition** | **Enterprise Edition** |
| Tabbed interface                  | + | + |
| Dark theme                        | + | + |
| 125%, 150%, 175%, 200% scaling    | + | + |
| White Label                       | - | - |
| Integrated test example (node.js) | + | + |
| Mobile web editors                | - | +* |
| **Plugins & Macros** | **Community Edition** | **Enterprise Edition** |
| Plugins                           | + | + |
| Macros                            | + | + |
| **Collaborative capabilities** | **Community Edition** | **Enterprise Edition** |
| Two co-editing modes              | + | + |
| Comments                          | + | + |
| Built-in chat                     | + | + |
| Review and tracking changes       | + | + |
| Display modes of tracking changes | + | + |
| Version history                   | + | + |
| **Document Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Adding Content control          | + | + |
| Editing Content control         | + | + |
| Layout tools                    | + | + |
| Table of contents               | + | + |
| Navigation panel                | + | + |
| Mail Merge                      | + | + |
| Comparing Documents             | + | + |
| **Spreadsheet Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Functions, formulas, equations  | + | + |
| Table templates                 | + | + |
| Pivot tables                    | + | + |
| Data validation                 | + | + |
| Conditional formatting          | + | + |
| Sparklines                      | + | + |
| Sheet Views                     | + | + |
| **Presentation Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Transitions                     | + | + |
| Animations                      | + | + |
| Presenter mode                  | + | + |
| Notes                           | + | + |
| **Form creator features** | **Community Edition** | **Enterprise Edition** |
| Adding form fields              | + | + |
| Form preview                    | + | + |
| Saving as PDF                   | + | + |
| **PDF Editor features**      | **Community Edition** | **Enterprise Edition** |
| Text editing and co-editing                                | + | + |
| Work with pages (adding, deleting, rotating)               | + | + |
| Inserting objects (shapes, images, hyperlinks, etc.)       | + | + |
| Text annotations (highlight, underline, cross out, stamps) | + | + |
| Comments                        | + | + |
| Freehand drawings               | + | + |
| Form filling                    | + | + |
| | [Get it now](https://www.onlyoffice.com/download-community?utm_source=github&utm_medium=cpc&utm_campaign=GitHubPipedrive#docs-community)  | [Start Free Trial](https://www.onlyoffice.com/download?utm_source=github&utm_medium=cpc&utm_campaign=GitHubPipedrive#docs-enterprise)  |

\* If supported by DMS.

## Need help? User Feedback and Support 💡

* **🐞 Found a bug?** Please report it by creating an [issue](https://github.com/ONLYOFFICE/onlyoffice-pipedrive/issues).
* **❓ Have a question?** Ask our community and developers on the [ONLYOFFICE Forum](https://community.onlyoffice.com).
* **👨‍💻 Need help for developers?** Check our [API documentation](https://api.onlyoffice.com).
* **💡 Want to suggest a feature?** Share your ideas on our [feedback platform](https://feedback.onlyoffice.com/forums/966080-your-voice-matters).