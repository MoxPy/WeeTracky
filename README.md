# WeeTracky


**WeeTracky is a robust backend, developed in Go, designed to streamline the supply chain management process for small manufacturing enterprises, particularly those involved in jewelry production. This application allows efficient tracking of products, materials, suppliers, and certifications, providing a comprehensive view of the entire supply chain.**

## Core Technologies

* **Go (Golang)**: The backend is written in Go, offering a high-performance and statically typed language for optimal efficiency.

* **MongoDB Atlas**: In an effort to simplify database management for small businesses, WeeTracky utilizes MongoDB Atlas. Tailored for ease of use and scalability, MongoDB Atlas offers a fully-managed cloud database service, eliminating the need for businesses to handle complex database administration tasks. This decision aims to empower small manufacturing enterprises with a user-friendly and efficient solution, allowing them to focus on their core operations while ensuring a robust and adaptable supply chain management system.

## Routes
#### Products

    /products/add: Add a new product to the system.
    /products/update: Update an existing product.
    /products/all: Retrieve a list of all products.
    /products/find-product: Find a specific product by ID.
    /products/find-by-material: Retrieve products based on the material used.
    /products/delete-product: Delete a product.

#### Materials

    /materials/add: Add a new material to the system.
    /materials/update: Update an existing material.
    /materials/all: Retrieve a list of all materials.
    /materials/find-material: Find a specific material by ID.
    /materials/find-by-supplier: Retrieve materials based on the supplier.
    /materials/delete-material: Delete a material.

#### Suppliers

    /suppliers/add: Add a new supplier to the system.
    /suppliers/update: Update an existing supplier.
    /suppliers/all: Retrieve a list of all suppliers.
    /suppliers/find-supplier: Find a specific supplier by ID.
    /suppliers/delete-supplier: Delete a supplier.

#### Certifications

    /certs/add: Add a new certification to the system.
    /certs/all: Retrieve a list of all certifications.

#### Company

    /company/init: Initialize the company within the system. Ensure the COMPANY variable is declared in the .env file.

## Models

Supply Track employs a well-structured set of models to represent key entities in the supply chain management system. Each model encapsulates specific attributes and relationships, enhancing the clarity and organization of the underlying data.

#### Product

- **ID**: Unique identifier for the product.
- **Name**: Name of the product.
- **MadeIn**: Manufacturing origin of the product.
- **Materials**: List of materials used in the product.
- **Price**: Price of the product.
- **Description**: Description of the product.
- **SustainablePackage**: Indicates whether the packaging is sustainable.

#### Material

- **ID**: Unique identifier for the material.
- **Name**: Name of the material.
- **Supplier**: Supplier information for the material.
- **Origin**: Origin information of the material.
- **Sustainable**: Indicates whether the material is sustainable.
- **Details**: Additional details about the material.
- **LastOrder**: Timestamp of the last order for the material.

#### Supplier

- **ID**: Unique identifier for the supplier.
- **Name**: Name of the supplier.
- **Country**: Country of the supplier.
- **City**: City of the supplier.

#### Cert

- **ID**: Unique identifier for the certification.
- **Name**: Name of the certification.
- **Issuer**: Entity issuing the certification.
- **Details**: Additional details about the certification.

#### Company

- **ID**: Unique identifier for the company.
- **Name**: Name of the company.
- **Products**: List of products associated with the company.
- **Materials**: List of materials associated with the company.
- **Suppliers**: List of suppliers associated with the company.
- **Certs**: List of certifications associated with the company.

The ID follows a specific format, starting with a designated letter assigned to the respective model.

## Usage

    git clone https://github.com/MoxPy/WeeTracky.git
    cd WeeTracky

### Launch the Application:

    go run main.go

    STEP 1: 
    Declare Environment Variables
    -
    Create a .env file with MongoDB connection details and declare the COMPANY variable.

    MONGODB_URI=mongodb+srv://your-username:your-password@your-cluster.mongodb.net/your-database?retryWrites=true
    DB_NAME=your-database
    COLLECTION_NAME=Products
    COMPANY=YourCompany

    STEP 2
    Initialize Your Company
    -
    Execute a POST request to localhost:8080/company/init to set up your company.
    Start adding your products, materials, suppliers and certifications via the APIs.

## Contribution

If you wish to contribute to this project, fork the repository, implement your changes, and submit a pull request. The project is licensed under the terms of the GNU Affero General Public License v3.0.

For more information about the license, please refer to the [LICENSE](LICENSE) file in the project.

## TODO:

- **Input Validation**: Enhance the validation of input data to ensure the reliability and security of the application.

- **Database Support**: Explore and implement support for additional databases, providing users with more flexibility.

- **Advanced Search Filters**: Extend search functionalities with additional filters, improving the precision of data retrieval.

- **Performance Optimization**: Evaluate and optimize the application's performance for an even smoother user experience.

- **Docker Support**: Docker support for easier deployment and consistent environments across different systems.

- **Authentication System**

## Disclaimer:
WeeTracky is provided as-is, without any warranties or guarantees of any kind, expressed or implied. The use of this application is at your own risk, and the developer disclaims any responsibility for any damages or losses that may arise from its use.

While efforts have been made to ensure the reliability and accuracy of the code, it is essential to review and test thoroughly before deploying in a production environment. The developer is not liable for any consequences, including but not limited to data loss, system failures, or other issues that may occur during the use of WeeTracky.

Users are encouraged to contribute to the project, report issues, and participate in discussions. However, the developer reserves the right to make changes to the project without prior notice.

**By using WeeTracky, you agree to these terms and acknowledge that the application is intended for local use. Any future updates or modifications, including the potential implementation of an authentication system, will be at the discretion of the developer.**

For questions, commercial inquiries or additional information, feel free to contact me via [LinkedIn](https://www.linkedin.com/in/manuel-lanzani-59071b251/).
